package service

import (
	"fmt"
	"sync"
	"time"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/pkg/db"
	"forestwhisperGo/apps/api/pkg/utils"
)

var (
	cozeClientCache *utils.CozeClient
	cozeMutex       sync.RWMutex
)

func GetCozeClient() *utils.CozeClient {
	cozeMutex.RLock()
	if cozeClientCache != nil {
		cozeMutex.RUnlock()
		return cozeClientCache
	}
	cozeMutex.RUnlock()

	cozeMutex.Lock()
	defer cozeMutex.Unlock()

	if cozeClientCache != nil {
		return cozeClientCache
	}

	accessToken := utils.GetString("coze", "access_token")
	botID := utils.GetString("coze", "bot_id")

	cozeClientCache = utils.NewCozeClient(accessToken, botID)
	return cozeClientCache
}

func ResetCozeClient() {
	cozeMutex.Lock()
	defer cozeMutex.Unlock()
	cozeClientCache = nil
}

func GetUserConversations(uid int) ([]models.AIConversation, error) {
	var conversations []models.AIConversation
	err := db.GetDB().Where("user_id = ?", uid).Order("updated_at desc").Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func GetConversationByID(uid int, conversationID string) *models.AIConversation {
	var conv models.AIConversation
	if err := db.GetDB().Where("user_id = ? AND conversation_id = ?", uid, conversationID).First(&conv).Error; err != nil {
		return nil
	}
	return &conv
}

func CreateConversation(uid int, conversationID, title string) error {
	now := int(time.Now().Unix())
	conv := models.AIConversation{
		UserId:         uid,
		ConversationId: conversationID,
		Title:          title,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return db.GetDB().Create(&conv).Error
}

func UpdateConversationTime(conversationID string) error {
	now := int(time.Now().Unix())
	result := db.GetDB().Model(&models.AIConversation{}).
		Where("conversation_id = ?", conversationID).
		Update("updated_at", now)
	return result.Error
}

func UpdateConversationTitle(conversationID, title string) error {
	result := db.GetDB().Model(&models.AIConversation{}).
		Where("conversation_id = ?", conversationID).
		Update("title", title)
	return result.Error
}

func DeleteConversation(uid int, conversationID string) error {
	result := db.GetDB().Where("user_id = ? AND conversation_id = ?", uid, conversationID).Delete(&models.AIConversation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("会话不存在")
	}
	return nil
}

func GenerateTitle(firstMessage string) string {
	runes := []rune(firstMessage)
	if len(runes) > 30 {
		return string(runes[:30]) + "..."
	}
	return firstMessage
}
