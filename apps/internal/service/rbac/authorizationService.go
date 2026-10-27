package rbac

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"forestwhisperGo/apps/internal/models"
	"forestwhisperGo/apps/internal/models/rbac"
	"forestwhisperGo/apps/pkg/constants"
	"forestwhisperGo/apps/pkg/db"
)

const (
	permCachePrefix = "rbac:perms:"
	roleCachePrefix = "rbac:roles:"
	superadminCache = "rbac:superadmin:"
	cacheTTL        = 5 * time.Minute
)

func IsSuperadmin(user *models.User) bool {
	cacheKey := superadminCache + strconv.Itoa(user.Uid)
	if cached, err := db.Get(cacheKey); err == nil && cached == "1" {
		return true
	}

	var count int64
	db.GetDB().Model(&rbac.UserRole{}).
		Where("user_id = ? AND role_name = ?", user.Uid, string(constants.RoleSuperadmin)).
		Count(&count)
	if count > 0 {
		db.SetEX(cacheKey, "1", int(cacheTTL.Seconds()))
		return true
	}
	return false
}

func GetUserRoles(userId int) []string {
	cacheKey := roleCachePrefix + strconv.Itoa(userId)
	if cached, err := db.Get(cacheKey); err == nil && cached != "" {
		var roles []string
		if json.Unmarshal([]byte(cached), &roles) == nil {
			return roles
		}
	}

	var rows []rbac.UserRole
	db.GetDB().Where("user_id = ?", userId).Find(&rows)
	roleNames := make([]string, 0, len(rows))
	for _, row := range rows {
		roleNames = append(roleNames, row.RoleName)
	}

	if data, err := json.Marshal(roleNames); err == nil {
		db.SetEX(cacheKey, string(data), int(cacheTTL.Seconds()))
	}
	return roleNames
}

func HasRole(user *models.User, roleName string) bool {
	for _, r := range GetUserRoles(user.Uid) {
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
	return HasRole(user, string(constants.RoleAdmin))
}

func HasPermission(user *models.User, permissionName string) bool {
	if IsSuperadmin(user) {
		return true
	}
	for _, p := range GetUserPermissions(user.Uid) {
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
	cacheKey := permCachePrefix + strconv.Itoa(userId)
	if cached, err := db.Get(cacheKey); err == nil && cached != "" {
		var perms []string
		if json.Unmarshal([]byte(cached), &perms) == nil {
			return perms
		}
	}

	seen := map[constants.Permission]struct{}{}
	for _, roleName := range GetUserRoles(userId) {
		for _, perm := range constants.RolePermissions[constants.Role(roleName)] {
			seen[perm] = struct{}{}
		}
	}
	perms := make([]string, 0, len(seen))
	for perm := range seen {
		perms = append(perms, string(perm))
	}
	sort.Strings(perms)

	if data, err := json.Marshal(perms); err == nil {
		db.SetEX(cacheKey, string(data), int(cacheTTL.Seconds()))
	}
	return perms
}

func InvalidateUserCache(userId int) {
	ctx := context.Background()
	client := db.Client()
	key := strconv.Itoa(userId)
	client.Del(ctx, permCachePrefix+key)
	client.Del(ctx, roleCachePrefix+key)
	client.Del(ctx, superadminCache+key)
}
