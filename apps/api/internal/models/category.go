package models

type Category struct {
	Mid         int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:150;not null"`
	Slug        string `gorm:"size:150;not null;uniqueIndex"`
	Description string `gorm:"size:150;default:''"`
	Parent      int    `gorm:"default:0"`
	Order       int    `gorm:"default:0"`
}

func (Category) TableName() string { return "categories" }
