package rbac

import (
	"errors"
	"fmt"
	"time"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/internal/models/rbac"
	"forestwhisperGo/apps/api/pkg/constants"
	"forestwhisperGo/apps/api/pkg/db"

	"gorm.io/gorm"
)

func AssignRoleToUser(userId int, roleName string, grantedBy int) bool {
	if !constants.IsValidRole(roleName) {
		return false
	}
	var user models.User
	if db.GetDB().First(&user, userId).Error != nil {
		return false
	}
	return assignRoleTx(db.GetDB(), userId, roleName, grantedBy) == nil
}

func AssignRoleToUserTx(tx *gorm.DB, userId int, roleName string, grantedBy int) error {
	if !constants.IsValidRole(roleName) {
		return fmt.Errorf("invalid role: %s", roleName)
	}
	return assignRoleTx(tx, userId, roleName, grantedBy)
}

func assignRoleTx(tx *gorm.DB, userId int, roleName string, grantedBy int) error {
	var existing rbac.UserRole
	err := tx.Where("user_id = ? AND role_name = ?", userId, roleName).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(&rbac.UserRole{
		UserId:    userId,
		RoleName:  roleName,
		GrantedAt: int(time.Now().Unix()),
		GrantedBy: grantedBy,
	}).Error
}

func RemoveRoleFromUser(userId int, roleName string) bool {
	db.GetDB().Where("user_id = ? AND role_name = ?", userId, roleName).Delete(&rbac.UserRole{})
	return true
}

func CheckUserPermission(userId int, permissionName string) bool {
	var user models.User
	if db.GetDB().First(&user, userId).Error != nil {
		return false
	}
	return HasPermission(&user, permissionName)
}
