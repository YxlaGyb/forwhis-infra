package models

type Relationship struct {
	Cid int `gorm:"primaryKey"`
	Mid int `gorm:"primaryKey"`
}

func (Relationship) TableName() string { return "relationship" }
