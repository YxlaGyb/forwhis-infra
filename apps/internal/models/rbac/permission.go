package rbac

type Permission struct {
	Id          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:64;unique;not null"`
	Description string `gorm:"type:text;default:''"`
	Category    string `gorm:"size:32;default:'默认'"`
	CreatedAt   int    `gorm:"default:0"`
}

func (Permission) TableName() string { return "permission" }
