package rbac

type RolePermission struct {
	Id           int `gorm:"primaryKey;autoIncrement"`
	RoleId       int `gorm:"not null;uniqueIndex:uix_role_permission"`
	PermissionId int `gorm:"not null;uniqueIndex:uix_role_permission"`
}

func (RolePermission) TableName() string { return "rolePermission" }
