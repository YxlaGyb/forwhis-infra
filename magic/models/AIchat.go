package models

type AIConversation struct {
	Id             int    `gorm:"primaryKey" json:"id"`
	UserId         int    `gorm:"index;not null" json:"userId"`
	ConversationId string `gorm:"size:100;index;not null" json:"conversationId"`
	Title          string `gorm:"size:200" json:"title"`
	CreatedAt      int    `gorm:"default:0" json:"createdAt"`
	UpdatedAt      int    `gorm:"default:0" json:"updatedAt"`
}

func (AIConversation) TableName() string { return "ai_conversation" }
