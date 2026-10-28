package models

type Userinfo struct {
	Uiuid   int    `gorm:"primaryKey;autoIncrement"`
	Uid     int    `gorm:"not null;index"`
	Type    string `gorm:"size:50;not null"`
	Content string `gorm:"default:''"`
	Time    int    `gorm:"default:0"`
}

func (Userinfo) TableName() string { return "userinfo" }
