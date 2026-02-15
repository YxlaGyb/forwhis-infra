package models

type User struct {
	Uid        int    `gorm:"primaryKey"`
	Name       string `gorm:"size:32"`
	Mail       string `gorm:"size:150"`
	Password   string `gorm:"size:255"`
	Url        string `gorm:"size:150"`
	CreatedAt  int    `gorm:"default:0"`
	LastLogin  int    `gorm:"default:0"`
	LastIp     string `gorm:"size:150;default:''"`
	IsActive   int    `gorm:"default:0"`
	IsLoggedIn int    `gorm:"default:0"`
}

func (User) TableName() string { return "user" }
