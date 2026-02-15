package models

type Meta struct {
	Mid         int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:150;not null"`
	Slug        string `gorm:"size:150;not null"`
	Type        string `gorm:"size:32;not null"`
	Description string `gorm:"size:150;default:''"`
	Count       int    `gorm:"default:0"`
	Order       int    `gorm:"default:0"`
	Parent      int    `gorm:"default:0"`
}

const (
	MetaTypeCategory = "category"
	MetaTypeTag      = "tag"
)

func (Meta) TableName() string { return "meta" }
