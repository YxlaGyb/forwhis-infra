package rbac

import (
	"context"
	"encoding/json"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/pkg/constants"
	"forestwhisperGo/pkg/db"
)

const (
	permCachePrefix = "rbac:perms:"
	roleCachePrefix = "rbac:roles:"
	superadminCache = "rbac:superadmin:"
	cacheTTL        = 5 * time.Minute
)

func IsSuperadmin(user *models.User) bool {
	cacheKey := superadminCache + intToString(user.Uid)
	if cached, err := db.Get(cacheKey); err == nil && cached == "1" {
		return true
	}

	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	for _, ur := range roles {
		if ur.Role.Name == string(constants.RoleSuperadmin) {
			db.SetEX(cacheKey, "1", int(cacheTTL.Seconds()))
			return true
		}
	}
	return false
}

func GetUserRoles(userId int) []string {
	cacheKey := roleCachePrefix + intToString(userId)
	if cached, err := db.Get(cacheKey); err == nil && cached != "" {
		var roles []string
		if json.Unmarshal([]byte(cached), &roles) == nil {
			return roles
		}
	}

	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", userId).Preload("Role").Find(&roles)
	roleNames := []string{}
	for _, ur := range roles {
		roleNames = append(roleNames, ur.Role.Name)
	}

	if data, err := json.Marshal(roleNames); err == nil {
		db.SetEX(cacheKey, string(data), int(cacheTTL.Seconds()))
	}
	return roleNames
}

func HasRole(user *models.User, roleName string) bool {
	roles := GetUserRoles(user.Uid)
	for _, r := range roles {
		if r == roleName {
			return true
		}
	}
	return false
}

func HasAnyAdminRole(user *models.User) bool {
	if IsSuperadmin(user) {
		return true
	}
	roles := GetUserRoles(user.Uid)
	for _, r := range roles {
		if r == string(constants.RoleAdmin) {
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
	cacheKey := permCachePrefix + intToString(userId)
	if cached, err := db.Get(cacheKey); err == nil && cached != "" {
		var perms []string
		if json.Unmarshal([]byte(cached), &perms) == nil {
			return perms
		}
	}

	var perms []string
	db.GetDB().Table("permission").
		Select("DISTINCT permission.name").
		Joins(`JOIN "rolePermission" ON "rolePermission".permission_id = permission.id`).
		Joins(`JOIN "userrole" ON "userrole".role_id = "rolePermission".role_id`).
		Where(`"userrole".user_id = ?`, userId).
		Pluck("name", &perms)

	if data, err := json.Marshal(perms); err == nil {
		db.SetEX(cacheKey, string(data), int(cacheTTL.Seconds()))
	}
	return perms
}

func InvalidateUserCache(userId int) {
	ctx := context.Background()
	client := db.Client()
	client.Del(ctx, permCachePrefix+intToString(userId))
	client.Del(ctx, roleCachePrefix+intToString(userId))
	client.Del(ctx, superadminCache+intToString(userId))
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	var negative bool
	if n < 0 {
		negative = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
