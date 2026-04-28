package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"forestwhisperGo/magic/interfaces/dto/request"
	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/pkg/db"
	"forestwhisperGo/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	MaxMessageLength   = 4000
	RateLimitMax       = 20
	RateLimitWindow    = time.Minute
	RateLimitKeyPrefix = "ai_chat_rate:"
)

func getClientIP(c *gin.Context) string {
	return c.ClientIP()
}

func checkRateLimit(key string) (bool, error) {
	redisKey := fmt.Sprintf("%s%s", RateLimitKeyPrefix, key)
	count, err := db.Incr(redisKey)
	if err != nil {
		return false, err
	}
	if count == 1 {
		db.Expire(context.Background(), redisKey, RateLimitWindow)
	}
	return count <= RateLimitMax, nil
}

func GetConversations(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONSuccess(c, "ok", []gin.H{})
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

	conversationID := c.Query("conversationId")
	if conversationID == "" {
		middleware.JSONError(c, "会话ID不能为空喵", nil)
		return
	}

	if user != nil {
		conv := service.GetConversationByID(user.Uid, conversationID)
		if conv == nil {
			middleware.JSONError(c, "会话不存在喵", nil)
			return
		}
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

	var rateLimitKey string
	if user != nil {
		rateLimitKey = fmt.Sprintf("%d", user.Uid)
	} else {
		rateLimitKey = getClientIP(c)
	}

	allowed, err := checkRateLimit(rateLimitKey)
	if err != nil {
		middleware.JSONError(c, "服务暂时不可用喵", nil)
		return
	}
	if !allowed {
		middleware.JSONError(c, "请求太频繁了, 请稍后再试喵", nil)
		return
	}

	var req request.ChatRequest

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
	if req.ConversationID != "" && user != nil {
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

	var userID string
	if user != nil {
		userID = fmt.Sprintf("user_%d", user.Uid)
	} else {
		userID = fmt.Sprintf("guest_%s", getClientIP(c))
	}
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

	if user == nil {
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
		middleware.JSONError(c, "请先登录喵", nil)
		return
	}

	var req request.DeleteConversationRequest

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
