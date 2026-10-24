package rbac

import (
	"time"

	"forestwhisperGo/apps/internal/models"
	"forestwhisperGo/apps/internal/models/rbac"
	"forestwhisperGo/apps/pkg/db"

	"gorm.io/gorm"
)

func GetOrCreateRole(roleName, description string) *rbac.Role {
	var role rbac.Role
	if err := db.GetDB().Where("name = ?", roleName).First(&role).Error; err != nil {
		now := int(time.Now().Unix())
		role = rbac.Role{Name: roleName, Description: description, CreatedAt: now, UpdatedAt: now}
		db.GetDB().Create(&role)
	}
	return &role
}

func GetOrCreatePermission(name, description, category string) *rbac.Permission {
	var perm rbac.Permission
	if err := db.GetDB().Where("name = ?", name).First(&perm).Error; err != nil {
		perm = rbac.Permission{Name: name, Description: description, Category: category, CreatedAt: int(time.Now().Unix())}
		db.GetDB().Create(&perm)
	}
	return &perm
}

func AssignRoleToUser(userId int, roleName string, grantedBy int) bool {
	var user models.User
	if db.GetDB().First(&user, userId).Error != nil {
		return false
	}
	var role rbac.Role
	if db.GetDB().Where("name = ?", roleName).First(&role).Error != nil {
		return false
	}
	var existing rbac.UserRole
	if db.GetDB().Where("user_id = ? AND role_id = ?", userId, role.Id).First(&existing).Error == nil {
		return true
	}
	userRole := rbac.UserRole{UserId: userId, RoleId: role.Id, GrantedAt: int(time.Now().Unix()), GrantedBy: grantedBy}
	db.GetDB().Create(&userRole)
	return true
}

func AssignRoleToUserTx(tx *gorm.DB, userId int, roleName string, grantedBy int) error {
	var role rbac.Role
	if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}
	var existing rbac.UserRole
	if tx.Where("user_id = ? AND role_id = ?", userId, role.Id).First(&existing).Error == nil {
		return nil
	}
	userRole := rbac.UserRole{UserId: userId, RoleId: role.Id, GrantedAt: int(time.Now().Unix()), GrantedBy: grantedBy}
	return tx.Create(&userRole).Error
}

func RemoveRoleFromUser(userId int, roleName string) bool {
	var role rbac.Role
	if db.GetDB().Where("name = ?", roleName).First(&role).Error != nil {
		return false
	}
	db.GetDB().Where("user_id = ? AND role_id = ?", userId, role.Id).Delete(&rbac.UserRole{})
	return true
}

func GrantPermissionToRole(roleName, permissionName string) bool {
	var role rbac.Role
	if db.GetDB().Where("name = ?", roleName).First(&role).Error != nil {
		return false
	}
	var perm rbac.Permission
	if db.GetDB().Where("name = ?", permissionName).First(&perm).Error != nil {
		return false
	}
	var existing rbac.RolePermission
	if db.GetDB().Where("role_id = ? AND permission_id = ?", role.Id, perm.Id).First(&existing).Error == nil {
		return true
	}
	rp := rbac.RolePermission{RoleId: role.Id, PermissionId: perm.Id}
	db.GetDB().Create(&rp)
	return true
}

func RevokePermissionFromRole(roleName, permissionName string) bool {
	var role rbac.Role
	if db.GetDB().Where("name = ?", roleName).First(&role).Error != nil {
		return false
	}
	var perm rbac.Permission
	if db.GetDB().Where("name = ?", permissionName).First(&perm).Error != nil {
		return false
	}
	db.GetDB().Where("role_id = ? AND permission_id = ?", role.Id, perm.Id).Delete(&rbac.RolePermission{})
	return true
}

func CheckUserPermission(userId int, permissionName string) bool {
	var user models.User
	if db.GetDB().First(&user, userId).Error != nil {
		return false
	}
	return HasPermission(&user, permissionName)
}
