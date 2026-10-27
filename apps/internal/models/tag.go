package models

type Tag struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:150;not null;uniqueIndex"`
	Slug      string `gorm:"size:150;not null;uniqueIndex"`
	CreatedAt int    `gorm:"default:0"`
}

func (Tag) TableName() string { return "tags" }

type ContentTag struct {
	Cid   int `gorm:"primaryKey;autoIncrement:false"`
	TagId int `gorm:"primaryKey;autoIncrement:false;index"`
}

func (ContentTag) TableName() string { return "content_tags" }
