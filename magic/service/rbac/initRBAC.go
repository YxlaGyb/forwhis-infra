package rbac

import (
	"forestwhisperGo/magic/constants"
)

func InitDefaultRbac() {
	for _, p := range constants.AllPermissions {
		GetOrCreatePermission(string(p.Name), p.Desc, string(p.Category))
	}
	for roleName, perms := range constants.RolePermissions {
		GetOrCreateRole(string(roleName), constants.RoleDescriptions[roleName])
		for _, perm := range perms {
			GrantPermissionToRole(string(roleName), string(perm))
		}
	}
}
