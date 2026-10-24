package rbac

type UserRole struct {
	Id        int  `gorm:"primaryKey;autoIncrement"`
	UserId    int  `gorm:"not null;uniqueIndex:uix_user_role"`
	RoleId    int  `gorm:"not null;uniqueIndex:uix_user_role"`
	GrantedAt int  `gorm:"default:0"`
	GrantedBy int  `gorm:"default:0"`
	Role      Role `gorm:"foreignKey:RoleId"`
}

func (UserRole) TableName() string { return "userrole" }
