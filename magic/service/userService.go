package service

import (
	"math/rand"
	"time"

	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"
)

func GenerateRandomName() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "小可爱" + string(b)
}

func GetUserByEmail(email string) *models.User {
	var user models.User
	if db.DB.Where("mail = ?", email).First(&user).Error != nil {
		return nil
	}
	return &user
}

func CreateUser(name, email, password, ip, role string, autoGenerateName bool) *models.User {
	if autoGenerateName {
		name = GenerateRandomName()
	}
	now := int(time.Now().Unix())
	user := models.User{
		Name: name, Mail: email, Password: password,
		CreatedAt: now, LastLogin: 0, LastIp: ip, IsActive: 1,
	}
	db.DB.Create(&user)
	rbacSvc.AssignRoleToUser(user.Uid, role, 0)
	createDefaultUsermeta(user.Uid)
	return &user
}

func createDefaultUsermeta(uid int) {
	now := int(time.Now().Unix())
	defaults := []struct{ key, val string }{
		{"gender", "0"}, {"birthday", ""}, {"bio", "系统默认签名，送给每一位小可爱~"},
	}
	for _, d := range defaults {
		db.DB.Create(&models.Usermeta{UserId: uid, MetaKey: d.key, MetaValue: d.val, CreatedAt: now, UpdatedAt: now})
	}
}

func CreateUserinfo(uid int, t, content string) {
	db.DB.Create(&models.Userinfo{Uid: uid, Type: t, Content: content, Time: int(time.Now().Unix())})
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
	db.DB.Model(user).Updates(map[string]interface{}{"last_login": now, "last_ip": ip})
	CreateUserinfo(user.Uid, "login", "登录成功, IP: "+ip)

	token, _ := utils.GenerateToken("", user.Uid, user.Mail)
	perms := rbacSvc.GetUserPermissions(user.Uid)
	var roles []rbac.UserRole
	db.DB.Where("user_id = ?", user.Uid).Preload("Role").Find(&roles)
	roleNames := []string{}
	for _, r := range roles {
		roleNames = append(roleNames, r.Role.Name)
	}
	return map[string]interface{}{
		"uid": user.Uid, "name": user.Name, "email": user.Mail,
		"roles": roleNames, "permissions": perms, "token": token,
	}, 0
}

func GetUsersList() []map[string]interface{} {
	var users []models.User
	db.DB.Find(&users)
	result := []map[string]interface{}{}
	for _, u := range users {
		var roles []rbac.UserRole
		db.DB.Where("user_id = ?", u.Uid).Preload("Role").Find(&roles)
		roleNames := []string{}
		for _, r := range roles {
			roleNames = append(roleNames, r.Role.Name)
		}
		result = append(result, map[string]interface{}{
			"uid": u.Uid, "name": u.Name, "email": u.Mail, "roles": roleNames, "createdAt": u.CreatedAt,
		})
	}
	return result
}

func UpdateUser(uid int, data map[string]interface{}) bool {
	user := &models.User{}
	if db.DB.First(user, uid).Error != nil {
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
	db.DB.Model(user).Updates(updates)
	return true
}

func DeleteUser(uid, selectedUid int) bool {
	user := &models.User{}
	if db.DB.First(user, selectedUid).Error != nil {
		return false
	}
	db.DB.Delete(user)
	CreateUserinfo(uid, "delUser", "删除了 "+user.Name+" 喵")
	return true
}

func ChangePassword(uid int, oldPwd, newPwd string) (bool, string) {
	user := &models.User{}
	if db.DB.First(user, uid).Error != nil {
		return false, "用户不存在喵喵"
	}
	if !utils.VerifyPassword(user.Password, oldPwd) {
		return false, "旧密码错误喵喵"
	}
	db.DB.Model(user).Update("password", utils.HashPassword(newPwd))
	CreateUserinfo(uid, "changePassword", "用户修改了密码喵")
	return true, ""
}

func GetUsermeta(uid int, metaKey string) string {
	var meta models.Usermeta
	if db.DB.Where("user_id = ? AND meta_key = ?", uid, metaKey).First(&meta).Error != nil {
		return ""
	}
	return meta.MetaValue
}

func GetUserProfile(uid int) map[string]string {
	var metas []models.Usermeta
	db.DB.Where("user_id = ?", uid).Find(&metas)
	profile := map[string]string{}
	for _, m := range metas {
		profile[m.MetaKey] = m.MetaValue
	}
	return profile
}
