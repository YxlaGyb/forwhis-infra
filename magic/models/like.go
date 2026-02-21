package models

type Like struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	UserID     int    `gorm:"not null;index"`
	TargetType string `gorm:"size:16;not null;index"` // 'content' 或 'comment'
	TargetID   int    `gorm:"not null;index"`
	Created    int    `gorm:"default:0"`
}

const (
	LikeTargetTypeContent = "content"
	LikeTargetTypeComment = "comment"
)

func (Like) TableName() string { return "likes" }
