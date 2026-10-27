package rbac

type UserRole struct {
	UserId    int    `gorm:"primaryKey;autoIncrement:false"`
	RoleName  string `gorm:"primaryKey;size:32"`
	GrantedAt int    `gorm:"default:0"`
	GrantedBy int    `gorm:"default:0"`
}

func (UserRole) TableName() string { return "user_roles" }
