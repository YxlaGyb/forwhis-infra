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
		if ur.Role.Name == string(constants.RoleSuperadmin) {
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

func HasPermissionConst(user *models.User, perm constants.Permission) bool {
	return HasPermission(user, string(perm))
}

func CanPerformAction(user *models.User, resource, action string) bool {
	return HasPermission(user, resource+":"+action)
}

func GetUserPermissions(userId int) []string {
	var perms []string
	db.GetDB().Table("permission").
		Select("DISTINCT permission.name").
		Joins(`JOIN "rolePermission" ON "rolePermission".permission_id = permission.id`).
		Joins(`JOIN "userrole" ON "userrole".role_id = "rolePermission".role_id`).
		Where(`"userrole".user_id = ?`, userId).
		Pluck("name", &perms)
	return perms
}
