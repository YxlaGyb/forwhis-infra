package rbac

import (
	"forestwhisperGo/magic/constants"
)

func InitDefaultRbac() {
	for _, p := range constants.AllPermissions {
		GetOrCreatePermission(p.Name, p.Desc, p.Category)
	}
	for roleName, perms := range constants.RolePermissions {
		GetOrCreateRole(roleName, constants.RoleDescriptions[roleName])
		for _, perm := range perms {
			GrantPermissionToRole(roleName, perm)
		}
	}
}
