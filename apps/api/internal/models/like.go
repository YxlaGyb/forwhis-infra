package models

type Like struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	UserID     int    `gorm:"not null;index;uniqueIndex:uix_like_user_target"`
	TargetType string `gorm:"size:16;not null;index;uniqueIndex:uix_like_user_target"`
	TargetID   int    `gorm:"not null;index;uniqueIndex:uix_like_user_target"`
	Created    int    `gorm:"default:0"`
}

const (
	LikeTargetTypeContent = "content"
	LikeTargetTypeComment = "comment"
)

func (Like) TableName() string { return "likes" }
