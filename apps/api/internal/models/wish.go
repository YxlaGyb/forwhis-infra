package models

type Wish struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Created  int    `gorm:"default:0;index:idx_wishes_feed,priority:2"`
	AuthorID int    `gorm:"index"`
	Author   string `gorm:"size:50;default:''"`
	Mail     string `gorm:"size:150;default:''"`
	Content  string `gorm:"type:text;not null"`
	Status   string `gorm:"size:16;default:'approved';index:idx_wishes_feed,priority:1"`
}

const (
	WishStatusApproved = "approved"
	WishStatusHidden   = "hidden"
)

func (Wish) TableName() string { return "wishes" }
