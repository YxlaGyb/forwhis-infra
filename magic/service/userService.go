package service

import (
	"strings"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"gorm.io/gorm"
)

type UserWithRoles struct {
	User  *models.User
	Roles []rbac.UserRole
}

func GenerateRandomName() string {
	return "小可爱" + utils.SecureRandomLowerNumeric(6)
}

func GetUserByEmail(email string) *models.User {
	var user models.User
	if db.GetDB().Where("mail = ?", email).First(&user).Error != nil {
		return nil
	}
	return &user
}

func GetUserWithRoles(uid int) (*UserWithRoles, bool) {
	var user models.User
	var roles []rbac.UserRole

	db.GetDB().Where("user_id = ?", uid).Preload("Role").Find(&roles)
	if db.GetDB().First(&user, uid).Error != nil || user.Uid == 0 {
		return nil, false
	}

	return &UserWithRoles{
		User:  &user,
		Roles: roles,
	}, true
}

func GetUsersByUsername(name, exactly string) []map[string]interface{} {
	var users []map[string]interface{}
	if exactly == "1" {
		user := GetUserByEmail(name)
		if user != nil {
			users = []map[string]interface{}{{
				"uid":   user.Uid,
				"name":  user.Name,
				"email": user.Mail,
			}}
		}
	} else {
		var userList []models.User
		db.GetDB().Where("name LIKE ?", "%"+escapeLike(name)+"%").Find(&userList)
		for _, u := range userList {
			users = append(users, map[string]interface{}{
				"uid":   u.Uid,
				"name":  u.Name,
				"email": u.Mail,
			})
		}
	}
	return users
}

func GetUserProfile(uid int) map[string]string {
	var metas []models.Usermeta
	db.GetDB().Where("user_id = ?", uid).Find(&metas)
	profile := map[string]string{}
	for _, m := range metas {
		profile[m.MetaKey] = m.MetaValue
	}
	return profile
}

func GetUsermeta(uid int, metaKey string) string {
	var meta models.Usermeta
	if db.GetDB().Where("user_id = ? AND meta_key = ?", uid, metaKey).First(&meta).Error != nil {
		return ""
	}
	return meta.MetaValue
}

func GetUsersList() []map[string]interface{} {
	var users []models.User
	db.GetDB().Find(&users)

	if len(users) == 0 {
		return []map[string]interface{}{}
	}

	userIDs := make([]int, len(users))
	for i, u := range users {
		userIDs[i] = u.Uid
	}

	var allRoles []rbac.UserRole
	db.GetDB().Where("user_id IN ?", userIDs).Preload("Role").Find(&allRoles)

	userRolesMap := make(map[int][]string)
	for _, ur := range allRoles {
		userRolesMap[ur.UserId] = append(userRolesMap[ur.UserId], ur.Role.Name)
	}

	result := []map[string]interface{}{}
	for _, u := range users {
		roleNames := userRolesMap[u.Uid]
		if roleNames == nil {
			roleNames = []string{}
		}
		result = append(result, map[string]interface{}{
			"uid": u.Uid, "name": u.Name, "email": u.Mail, "roles": roleNames, "createdAt": u.CreatedAt,
		})
	}
	return result
}

func CreateUser(name, email, password, ip, role string, autoGenerateName bool) (*models.User, error) {
	if autoGenerateName {
		name = GenerateRandomName()
	}

	now := int(time.Now().Unix())
	user := models.User{
		Name: name, Mail: email, Password: password,
		CreatedAt: now, LastLogin: 0, LastIp: ip, IsActive: 1,
	}
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := rbacSvc.AssignRoleToUserTx(tx, user.Uid, role, 0); err != nil {
			return err
		}
		return createDefaultUsermetaTx(tx, user.Uid)
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func createDefaultUsermetaTx(tx *gorm.DB, uid int) error {
	now := int(time.Now().Unix())
	defaults := []struct{ key, val string }{
		{"gender", "0"}, {"birthday", ""}, {"bio", "系统默认签名，送给每一位小可爱~"},
	}
	for _, d := range defaults {
		if err := tx.Create(&models.Usermeta{UserId: uid, MetaKey: d.key, MetaValue: d.val, CreatedAt: now, UpdatedAt: now}).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateUserinfo(uid int, t, content string) {
	db.GetDB().Create(&models.Userinfo{Uid: uid, Type: t, Content: content, Time: int(time.Now().Unix())})
}

func LoginUser(email, password, ip string) (map[string]interface{}, int) {
	user := GetUserByEmail(email)
	if user == nil {
		return nil, 10101
	}
	if !utils.VerifyPassword(user.Password, password) {
		return nil, 10102
	}
	now := int(time.Now().Unix())
	db.GetDB().Model(user).Updates(map[string]interface{}{"last_login": now, "last_ip": ip})
	CreateUserinfo(user.Uid, "login", "登录成功, IP: "+ip)

	token, err := utils.GenerateToken("", user.Uid, user.Mail)
	if err != nil {
		return nil, 10103
	}
	perms := rbacSvc.GetUserPermissions(user.Uid)
	var roles []rbac.UserRole
	db.GetDB().Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	roleNames := []string{}
	for _, r := range roles {
		roleNames = append(roleNames, r.Role.Name)
	}
	return map[string]interface{}{
		"uid": user.Uid, "name": user.Name, "email": user.Mail,
		"roles": roleNames, "permissions": perms, "token": token,
	}, 0
}

func UpdateUser(uid int, data map[string]interface{}) bool {
	user := &models.User{}
	if db.GetDB().First(user, uid).Error != nil {
		return false
	}
	updates := map[string]interface{}{}
	if v, ok := data["username"].(string); ok {
		updates["name"] = v
	}
	if v, ok := data["email"].(string); ok {
		updates["mail"] = v
	}
	if v, ok := data["password"].(string); ok {
		updates["password"] = v
	}
	db.GetDB().Model(user).Updates(updates)
	return true
}

func DeleteUser(uid, selectedUid int) bool {
	user := &models.User{}
	if db.GetDB().First(user, selectedUid).Error != nil {
		return false
	}
	db.GetDB().Delete(user)
	CreateUserinfo(uid, "delUser", "删除了 "+user.Name+" 喵")
	return true
}

func ChangePassword(uid int, oldPwd, newPwd string) (bool, string) {
	user := &models.User{}
	if db.GetDB().First(user, uid).Error != nil {
		return false, "用户不存在喵喵"
	}
	if !utils.VerifyPassword(user.Password, oldPwd) {
		return false, "旧密码错误喵喵"
	}
	hashedPwd, err := utils.HashPassword(newPwd)
	if err != nil {
		return false, "密码加密失败喵喵"
	}
	db.GetDB().Model(user).Update("password", hashedPwd)
	CreateUserinfo(uid, "changePassword", "用户修改了密码喵")
	return true, ""
}

func UpdateUsermeta(uid int, metaKey, metaValue string) (bool, string) {
	var meta models.Usermeta
	result := db.GetDB().Where("user_id = ? AND meta_key = ?", uid, metaKey).First(&meta)
	now := int(time.Now().Unix())
	if result.Error != nil {
		meta = models.Usermeta{
			UserId:    uid,
			MetaKey:   metaKey,
			MetaValue: metaValue,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := db.GetDB().Create(&meta).Error; err != nil {
			return false, "创建元数据失败喵"
		}
		return true, ""
	}
	meta.MetaValue = metaValue
	meta.UpdatedAt = now
	if err := db.GetDB().Save(&meta).Error; err != nil {
		return false, "更新元数据失败喵"
	}
	return true, ""
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
