package rbac

type Role struct {
	Id          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:32;unique;not null"`
	Description string `gorm:"type:text;default:''"`
	CreatedAt   int    `gorm:"default:0"`
	UpdatedAt   int    `gorm:"default:0"`
}

func (Role) TableName() string { return "role" }
