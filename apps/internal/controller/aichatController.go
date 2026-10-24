package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/request"
	"forestwhisperGo/apps/internal/service"
	"forestwhisperGo/apps/pkg/db"
	"forestwhisperGo/apps/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

const (
	MaxMessageLength   = 4000
	RateLimitMax       = 20
	RateLimitWindow    = time.Minute
	RateLimitKeyPrefix = "ai_chat_rate:"
)

func getClientIP(c fiber.Ctx) string {
	return middleware.ClientIP(c)
}

func writeSSE(w *bufio.Writer, data fiber.Map) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: message\ndata: %s\n\n", encoded); err != nil {
		return err
	}
	return w.Flush()
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

func GetConversations(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return middleware.JSONSuccess(c, "ok", []fiber.Map{})
	}

	conversations, err := service.GetUserConversations(user.Uid)
	if err != nil {
		return middleware.JSONError(c, "获取会话列表失败喵", nil)
	}
	return middleware.JSONSuccess(c, "ok", conversations)
}

func GetMessages(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)

	conversationID := c.Query("conversationId")
	if conversationID == "" {
		return middleware.JSONError(c, "会话ID不能为空喵", nil)
	}

	if user != nil {
		conv := service.GetConversationByID(user.Uid, conversationID)
		if conv == nil {
			return middleware.JSONError(c, "会话不存在喵", nil)
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
		return middleware.JSONError(c, "获取消息失败喵: "+err.Error(), nil)
	}

	result := []fiber.Map{}
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Type == "answer" || msg.Type == "question" {
			result = append(result, fiber.Map{
				"role":      msg.Role,
				"content":   msg.Content,
				"createdAt": msg.CreatedAt,
			})
		}
	}
	return middleware.JSONSuccess(c, "ok", result)
}

func ChatStream(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)

	var rateLimitKey string
	if user != nil {
		rateLimitKey = fmt.Sprintf("%d", user.Uid)
	} else {
		rateLimitKey = getClientIP(c)
	}

	allowed, err := checkRateLimit(rateLimitKey)
	if err != nil {
		return middleware.JSONError(c, "服务暂时不可用喵", nil)
	}
	if !allowed {
		return middleware.JSONError(c, "请求太频繁了, 请稍后再试喵", nil)
	}

	var req request.ChatRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if strings.TrimSpace(req.Message) == "" {
		return middleware.JSONError(c, "消息不能为空喵", nil)
	}
	if len(req.Message) > MaxMessageLength {
		return middleware.JSONError(c, fmt.Sprintf("消息长度不能超过%d字符喵", MaxMessageLength), nil)
	}
	if req.ConversationID != "" && user != nil {
		conv := service.GetConversationByID(user.Uid, req.ConversationID)
		if conv == nil {
			return middleware.JSONError(c, "会话不存在喵", nil)
		}
	}

	cozeClient := service.GetCozeClient()
	if cozeClient.AccessToken == "" || cozeClient.BotID == "" {
		return middleware.JSONError(c, "AI服务未配置喵", nil)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	var userID string
	if user != nil {
		userID = fmt.Sprintf("user_%d", user.Uid)
	} else {
		userID = fmt.Sprintf("guest_%s", getClientIP(c))
	}
	isNewConversation := req.ConversationID == ""
	savedConversationID := req.ConversationID

	ctx, cancel := context.WithTimeout(c.RequestCtx(), 180*time.Second)
	defer cancel()

	streamReq := utils.ChatStreamRequest{
		Context:        ctx,
		ConversationID: req.ConversationID,
		UserID:         userID,
		Message:        req.Message,
	}

	var writeErr error
	err = c.SendStreamWriter(func(w *bufio.Writer) {
		streamErr := cozeClient.ChatStream(streamReq, func(event string, data string) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			var message fiber.Map
			switch event {
			case "conversation.message.delta":
				content := utils.ParseMessageContent(data)
				if content == "" {
					return nil
				}
				message = fiber.Map{"type": "content", "content": content}
			case "conversation.message.completed":
				message = fiber.Map{"type": "done"}
			case "conversation.chat.created":
				convID := utils.ParseConversationID(data)
				if convID == "" {
					return nil
				}
				savedConversationID = convID
				message = fiber.Map{"type": "conversation", "conversationId": convID}
			case "error":
				message = fiber.Map{"type": "error", "error": utils.ParseError(data)}
			default:
				return nil
			}
			writeErr = writeSSE(w, message)
			return writeErr
		})

		if writeErr == nil && streamErr != nil && streamErr != context.Canceled {
			writeErr = writeSSE(w, fiber.Map{"type": "error", "error": streamErr.Error()})
		}

		if writeErr == nil && streamErr == nil && user != nil && savedConversationID != "" {
			if isNewConversation {
				if err := service.CreateConversation(user.Uid, savedConversationID, service.GenerateTitle(req.Message)); err != nil {
					log.Printf("save AI conversation %s: %v", savedConversationID, err)
				}
			} else {
				if err := service.UpdateConversationTime(savedConversationID); err != nil {
					log.Printf("update AI conversation %s: %v", savedConversationID, err)
				}
			}
		}
	})
	if err != nil {
		return err
	}
	return writeErr
}

func DeleteConversation(c fiber.Ctx) error {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		return middleware.JSONError(c, "请先登录喵", nil)
	}

	var req request.DeleteConversationRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if req.ConversationID == "" {
		return middleware.JSONError(c, "会话ID不能为空喵", nil)
	}

	err := service.DeleteConversation(user.Uid, req.ConversationID)
	if err != nil {
		return middleware.JSONError(c, err.Error(), nil)
	}
	return middleware.JSONSuccess(c, "删除成功喵", nil)
}
