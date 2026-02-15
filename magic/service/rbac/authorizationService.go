package rbac

import (
	"forestwhisperGo/magic/constants"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/magic/utils/db"
)

func IsSuperadmin(user *models.User) bool {
	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	for _, ur := range roles {
		if ur.Role.Name == constants.RoleSuperadmin {
			return true
		}
	}
	return false
}

func HasRole(user *models.User, roleName string) bool {
	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	for _, ur := range roles {
		if ur.Role.Name == roleName {
			return true
		}
	}
	return false
}

func HasPermission(user *models.User, permissionName string) bool {
	if IsSuperadmin(user) {
		return true
	}
	perms := GetUserPermissions(user.Uid)
	for _, p := range perms {
		if p == permissionName {
			return true
		}
	}
	return false
}

func CanPerformAction(user *models.User, resource, action string) bool {
	return HasPermission(user, resource+":"+action)
}

func GetUserPermissions(userId int) []string {
	var perms []string
	var userRoles []rbac.UserRole
	db.GetDB().Where("user_id = ?", userId).Find(&userRoles)
	for _, ur := range userRoles {
		var rolePerms []rbac.RolePermission
		db.GetDB().Where("role_id = ?", ur.RoleId).Find(&rolePerms)
		for _, rp := range rolePerms {
			var perm rbac.Permission
			db.GetDB().First(&perm, rp.PermissionId)
			perms = append(perms, perm.Name)
		}
	}
	return perms
}
