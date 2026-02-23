package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

const (
	MaxMessageLength   = 4000
	RateLimitMax       = 20
	RateLimitWindow    = time.Minute
	RateLimitKeyPrefix = "ai_chat_rate:"
)

func checkRateLimit(uid int) (bool, error) {
	key := fmt.Sprintf("%s%d", RateLimitKeyPrefix, uid)
	count, err := db.Incr(key)
	if err != nil {
		return false, err
	}
	if count == 1 {
		db.Expire(context.Background(), key, RateLimitWindow)
	}
	return count <= RateLimitMax, nil
}

func GetConversations(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵", nil)
		return
	}

	conversations, err := service.GetUserConversations(user.Uid)
	if err != nil {
		middleware.JSONError(c, "获取会话列表失败喵", nil)
		return
	}
	middleware.JSONSuccess(c, "ok", conversations)
}

func GetMessages(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵", nil)
		return
	}

	conversationID := c.Query("conversationId")
	if conversationID == "" {
		middleware.JSONError(c, "会话ID不能为空喵", nil)
		return
	}

	conv := service.GetConversationByID(user.Uid, conversationID)
	if conv == nil {
		middleware.JSONError(c, "会话不存在喵", nil)
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		var limitInt int
		if _, err := fmt.Sscanf(l, "%d", &limitInt); err == nil && limitInt > 0 && limitInt <= 100 {
			limit = limitInt
		}
	}

	cozeClient := service.GetCozeClient()
	messages, err := cozeClient.GetMessages(conversationID, &utils.GetMessagesOptions{Limit: limit})
	if err != nil {
		middleware.JSONError(c, "获取消息失败喵: "+err.Error(), nil)
		return
	}

	result := []gin.H{}
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Type == "answer" || msg.Type == "question" {
			result = append(result, gin.H{
				"role":      msg.Role,
				"content":   msg.Content,
				"createdAt": msg.CreatedAt,
			})
		}
	}

	middleware.JSONSuccess(c, "ok", result)
}

func ChatStream(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵", nil)
		return
	}

	allowed, err := checkRateLimit(user.Uid)
	if err != nil {
		middleware.JSONError(c, "服务暂时不可用喵", nil)
		return
	}
	if !allowed {
		middleware.JSONError(c, "请求太频繁了，请稍后再试喵", nil)
		return
	}

	var req types.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		middleware.JSONError(c, "消息不能为空喵", nil)
		return
	}
	if len(req.Message) > MaxMessageLength {
		middleware.JSONError(c, fmt.Sprintf("消息长度不能超过%d字符喵", MaxMessageLength), nil)
		return
	}
	if req.ConversationID != "" {
		conv := service.GetConversationByID(user.Uid, req.ConversationID)
		if conv == nil {
			middleware.JSONError(c, "会话不存在喵", nil)
			return
		}
	}

	cozeClient := service.GetCozeClient()
	if cozeClient.AccessToken == "" || cozeClient.BotID == "" {
		middleware.JSONError(c, "AI服务未配置喵", nil)
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userID := fmt.Sprintf("user_%d", user.Uid)
	isNewConversation := req.ConversationID == ""
	savedConversationID := req.ConversationID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 180*time.Second)
	defer cancel()

	streamReq := utils.ChatStreamRequest{
		Context:        ctx,
		ConversationID: req.ConversationID,
		UserID:         userID,
		Message:        req.Message,
	}

	streamErr := cozeClient.ChatStream(streamReq, func(event string, data string) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		switch event {
		case "conversation.message.delta":
			content := utils.ParseMessageContent(data)
			if content != "" {
				c.SSEvent("message", gin.H{
					"type":    "content",
					"content": content,
				})
				c.Writer.Flush()
			}
		case "conversation.message.completed":
			c.SSEvent("message", gin.H{
				"type": "done",
			})
			c.Writer.Flush()
		case "conversation.chat.created":
			convID := utils.ParseConversationID(data)
			if convID != "" {
				savedConversationID = convID
				c.SSEvent("message", gin.H{
					"type":           "conversation",
					"conversationId": convID,
				})
				c.Writer.Flush()
			}
		case "error":
			errMsg := utils.ParseError(data)
			c.SSEvent("message", gin.H{
				"type":  "error",
				"error": errMsg,
			})
			c.Writer.Flush()
		}
		return nil
	})

	if streamErr != nil {
		if streamErr != context.Canceled {
			c.SSEvent("message", gin.H{
				"type":  "error",
				"error": streamErr.Error(),
			})
			c.Writer.Flush()
		}
		return
	}
	if isNewConversation && savedConversationID != "" {
		title := service.GenerateTitle(req.Message)
		if err := service.CreateConversation(user.Uid, savedConversationID, title); err != nil {
			// log error silently
		}
	} else if savedConversationID != "" {
		if err := service.UpdateConversationTime(savedConversationID); err != nil {
			// log error silently
		}
	}
}

func DeleteConversation(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵", nil)
		return
	}

	var req types.DeleteConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}
	if req.ConversationID == "" {
		middleware.JSONError(c, "会话ID不能为空喵", nil)
		return
	}

	err := service.DeleteConversation(user.Uid, req.ConversationID)
	if err != nil {
		middleware.JSONError(c, err.Error(), nil)
		return
	}

	middleware.JSONSuccess(c, "删除成功喵", nil)
}
