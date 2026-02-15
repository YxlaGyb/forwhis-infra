package controller

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/rbac"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

const codeExpirationTime = 300

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	clientIP := c.ClientIP()
	if !utils.IsValidEmail(req.Email) || !utils.IsValidPassword(req.Password) {
		middleware.JSONError(c, "邮箱或密码格式不正确喵", nil)
		return
	}

	cdKey := db.FormatKey("login:cd", req.Email)
	if exists, _ := db.Exists(cdKey); exists > 0 {
		middleware.JSONError(c, "操作过于频繁，请稍后再试喵", nil)
		return
	}
	db.SetEX(cdKey, "1", 60)

	result, errCode := service.LoginUser(req.Email, req.Password, clientIP)
	if errCode != 0 {
		db.Del(cdKey)
		middleware.JSONError(c, "邮箱或密码错误喵", nil)
		return
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	middleware.JSONSuccess(c, "登录成功喵", gin.H{
		"uid": result["uid"], "name": result["name"], "email": result["email"],
	})
}

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Code     string `json:"code"`
		CodeSalt string `json:"codeSalt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if !utils.IsValidEmail(req.Email) || len(req.Username) < 2 || len(req.Username) > 50 || len(req.Password) < 8 {
		middleware.JSONError(c, "用户名或密码格式不正确喵", nil)
		return
	}
	if len(req.Code) != 6 {
		middleware.JSONError(c, "验证码应为6位字母+数字喵喵", nil)
		return
	}
	if service.GetUserByEmail(req.Email) != nil {
		middleware.JSONError(c, "该邮箱已被注册喵喵", nil)
		return
	}

	valid, errMsg := verifyCode(req.Email, req.Code, req.CodeSalt)
	if !valid {
		middleware.JSONError(c, errMsg, nil)
		return
	}

	passwordHash := utils.HashPassword(req.Password)
	clientIP := c.ClientIP()

	service.CreateUser(req.Username, req.Email, passwordHash, clientIP, "user", true)

	result, errCode := service.LoginUser(req.Email, req.Password, clientIP)
	if errCode != 0 {
		middleware.JSONError(c, "注册后登录失败喵", nil)
		return
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	middleware.JSONSuccess(c, "注册成功喵", nil)
}

func SendEmailCodeRegister(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if !utils.IsValidEmail(req.Email) {
		middleware.JSONError(c, "邮箱格式错误喵", nil)
		return
	}
	if service.GetUserByEmail(req.Email) != nil {
		middleware.JSONError(c, "您的邮箱已经被使用了喵, 请换一个试试喵", nil)
		return
	}

	code := generateCode(6)
	codeSalt := utils.HashPassword(code)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; padding: 20px; border: 1px solid #eee; border-radius: 8px;">
			<h2 style="color: #333;">✨ 注册验证码待查收</h2>
			<p>您的验证码是 <b style="font-size: 24px; color: #007bff;">%s</b></p>
			<p style="color: #666;">请在 %d 分钟内完成验证，请勿泄露喵~</p>
			<hr style="border: none; border-top: 1px solid #eee;">
			<footer style="font-size: 12px; color: #999;">来自：数数洞洞平台</footer>
		</div>
	`, code, codeExpirationTime/60)

	if !utils.SendMail("注册验证码", []string{req.Email}, htmlContent) {
		middleware.JSONSystemError(c, "邮件服务连接超时, 请稍后再试喵喵", nil)
		return
	}

	codeKey := db.FormatKey("verify:code", req.Email)
	codeData := fmt.Sprintf("%s:%s", code, codeSalt)
	db.SetEX(codeKey, codeData, codeExpirationTime)

	middleware.JSONSuccess(c, "验证码已发送喵", gin.H{"codeSalt": codeSalt})
}

func GetUserByUsername(c *gin.Context) {
	name := c.Query("name")
	exactly := c.Query("exactly")

	if strings.TrimSpace(name) == "" {
		middleware.JSONError(c, "用户名不能为空喵", nil)
		return
	}

	var users []map[string]interface{}
	if exactly == "1" {
		user := service.GetUserByEmail(name)
		if user != nil {
			users = []map[string]interface{}{{
				"uid":   user.Uid,
				"name":  user.Name,
				"email": user.Mail,
			}}
		}
	} else {
		var userList []models.User
		db.GetDB().Where("name LIKE ?", "%"+name+"%").Find(&userList)
		for _, u := range userList {
			users = append(users, map[string]interface{}{
				"uid":   u.Uid,
				"name":  u.Name,
				"email": u.Mail,
			})
		}
	}

	middleware.JSONSuccess(c, "获取成功喵", users)
}

func DeleteUserByUid(c *gin.Context) {
	var req struct {
		Uid int `json:"uid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if req.Uid == 0 {
		middleware.JSONError(c, "用户ID不能为空喵", nil)
		return
	}

	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		middleware.JSONError(c, "未登录或登录已过期喵喵", nil)
		return
	}
	if currentUser.Uid == req.Uid {
		middleware.JSONError(c, "不能删除自己喵", nil)
		return
	}

	var targetUser models.User
	var targetRoles []rbac.UserRole
	db.GetDB().Where("user_id = ?", req.Uid).Preload("Role").Find(&targetRoles)
	db.GetDB().First(&targetUser, req.Uid)
	if targetUser.Uid == 0 {
		middleware.JSONError(c, "用户不存在喵喵", nil)
		return
	}

	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)
	isTargetSuperadmin := rbacSvc.IsSuperadmin(&targetUser)
	if isTargetSuperadmin && !isCurrentSuperadmin {
		middleware.JSONError(c, "无权删除超级管理员喵", nil)
		return
	}

	success := service.DeleteUser(currentUser.Uid, req.Uid)
	if !success {
		middleware.JSONError(c, "用户删除失败喵喵", nil)
		return
	}

	middleware.JSONSuccess(c, "用户删除成功喵", nil)
}

func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		middleware.JSONError(c, "旧密码和新密码不能为空喵", nil)
		return
	}
	if len(req.NewPassword) < 8 {
		middleware.JSONError(c, "新密码长度应不少于8个字符喵喵", nil)
		return
	}

	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵喵", nil)
		return
	}

	success, errMsg := service.ChangePassword(user.Uid, req.OldPassword, req.NewPassword)
	if !success {
		middleware.JSONError(c, errMsg, nil)
		return
	}

	middleware.JSONSuccess(c, "密码修改成功喵", nil)
}

func UpdateProfile(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	user := middleware.GetCurrentUser(c)
	if user == nil {
		middleware.JSONError(c, "未登录或登录已过期喵喵", nil)
		return
	}

	allowedKeys := map[string]int{
		models.MetaGender:   0,
		models.MetaBirthday: 0,
		models.MetaBio:      200,
	}

	for key, maxLen := range allowedKeys {
		if value, ok := req[key]; ok {
			if maxLen > 0 && len(value) > maxLen {
				middleware.JSONError(c, key+"长度不能超过"+strconv.Itoa(maxLen)+"个字符喵", nil)
				return
			}
			if key == models.MetaGender {
				gender, _ := strconv.Atoi(value)
				if gender < 0 || gender > 3 {
					middleware.JSONError(c, "性别值无效喵", nil)
					return
				}
			}
			success, errMsg := service.UpdateUsermeta(user.Uid, key, value)
			if !success {
				middleware.JSONError(c, errMsg, nil)
				return
			}
		}
	}

	middleware.JSONSuccess(c, "资料已提交审核或更新喵", nil)
}

func generateCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func verifyCode(email, code, codeSalt string) (bool, string) {
	codeKey := db.FormatKey("verify:code", email)
	storedData, err := db.Get(codeKey)
	if err != nil || storedData == "" {
		return false, "验证码不存在或已过期喵喵"
	}

	parts := strings.Split(storedData, ":")
	if len(parts) != 2 {
		return false, "验证码格式错误喵喵"
	}

	storedCode, storedSalt := parts[0], parts[1]
	if code != storedCode || codeSalt != storedSalt {
		return false, "验证码错误喵喵"
	}

	db.Del(codeKey)
	return true, ""
}
