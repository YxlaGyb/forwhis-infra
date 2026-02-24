package controller

import (
	"net/http"
	"strconv"
	"strings"

	"forestwhisperGo/magic/interfaces/http/middleware"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	rbacSvc "forestwhisperGo/magic/service/rbac"
	"forestwhisperGo/pkg/db"
	"forestwhisperGo/pkg/utils"

	"github.com/gin-gonic/gin"
)

const codeExpirationTime = 300

func Login(c *gin.Context) {
	var req types.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	clientIP := c.ClientIP()
	if !utils.IsValidEmail(req.Email) || !utils.IsValidPassword(req.Passwd) {
		middleware.JSONError(c, "邮箱或密码格式不正确喵", nil)
		return
	}

	cdKey := db.FormatKey("login:cd", req.Email)
	if exists, _ := db.Exists(cdKey); exists > 0 {
		middleware.JSONError(c, "操作过于频繁，请稍后再试喵", nil)
		return
	}
	db.SetEX(cdKey, "1", 60)

	result, errCode := service.LoginUser(req.Email, req.Passwd, clientIP)
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

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "forestwhisper",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	middleware.JSONSuccess(c, "退出成功喵", nil)
}

func Register(c *gin.Context) {
	if !service.GetOptionBool("feature_enableRegistration", true) {
		middleware.JSONError(c, "当前系统已关闭用户注册功能喵", nil)
		return
	}

	var req types.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}
	if !utils.IsValidEmail(req.Email) || !utils.IsValidName(req.Username) || !utils.IsValidPassword(req.Passwd) {
		middleware.JSONError(c, "邮箱、用户名或密码格式不正确喵（密码需8-72位，包含大小写字母和数字）", nil)
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

	valid, errMsg := verifyCode(req.Email, req.Code)
	if !valid {
		middleware.JSONError(c, errMsg, nil)
		return
	}

	clientIP := c.ClientIP()
	hashedPwd, err := utils.HashPassword(req.Passwd)
	if err != nil {
		middleware.JSONSystemError(c, "密码加密失败喵喵", nil)
		return
	}
	_, err = service.CreateUser(
		req.Username,
		req.Email,
		hashedPwd,
		clientIP, "user",
		true)

	if err != nil {
		middleware.JSONSystemError(c, "注册失败喵喵", nil)
		return
	}

	result, errCode := service.LoginUser(req.Email, req.Passwd, clientIP)
	if errCode != 0 {
		middleware.JSONError(c, "注册后登录失败喵", nil)
		return
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	middleware.JSONSuccess(c, "注册成功喵", nil)
}

func SendEmailCodeRegister(c *gin.Context) {
	if !service.GetOptionBool("feature_enableRegistration", true) {
		middleware.JSONError(c, "当前系统已关闭用户注册功能喵", nil)
		return
	}

	var req types.SendEmailCodeRequest

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

	code := utils.SecureRandomString(6)
	hashedCode, err := utils.HashPassword(code)
	if err != nil {
		middleware.JSONSystemError(c, "验证码生成失败喵喵", nil)
		return
	}
	htmlContent := `
		<div style="font-family: sans-serif; padding: 20px; border: 1px solid #eee; border-radius: 8px;">
			<h2 style="color: #333;">✨ 注册验证码待查收</h2>
			<p>您的验证码是 <b style="font-size: 24px; color: #007bff;">` + code + `</b></p>
			<p style="color: #666;">请在 5 分钟内完成验证，请勿泄露喵~</p>
			<hr style="border: none; border-top: 1px solid #eee;">
			<footer style="font-size: 12px; color: #999;">来自：数数洞洞平台</footer>
		</div>
	`

	if !utils.SendMail("注册验证码", []string{req.Email}, htmlContent) {
		middleware.JSONSystemError(c, "邮件服务连接超时, 请稍后再试喵喵", nil)
		return
	}

	codeKey := db.FormatKey("verify:code", req.Email)
	db.SetEX(codeKey, hashedCode, codeExpirationTime)
	middleware.JSONSuccess(c, "验证码已发送喵", gin.H{"codeSalt": hashedCode})
}

func GetUserByUsername(c *gin.Context) {
	name := c.Query("name")
	exactly := c.Query("exactly")

	if strings.TrimSpace(name) == "" {
		middleware.JSONError(c, "用户名不能为空喵", nil)
		return
	}

	users := service.GetUsersByUsername(name, exactly)
	middleware.JSONSuccess(c, "ok", users)
}

func DeleteUserByUid(c *gin.Context) {
	var req types.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if req.Uid == 0 {
		middleware.JSONError(c, "用户ID不能为空喵", nil)
		return
	}

	currentUser := middleware.GetCurrentUserFromContext(c)
	if currentUser.Uid == req.Uid {
		middleware.JSONError(c, "不能删除自己喵", nil)
		return
	}

	userWithRoles, ok := service.GetUserWithRoles(req.Uid)
	if !ok {
		middleware.JSONError(c, "用户不存在喵喵", nil)
		return
	}

	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)
	isTargetSuperadmin := rbacSvc.IsSuperadmin(userWithRoles.User)
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
	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		middleware.JSONError(c, "旧密码和新密码不能为空喵", nil)
		return
	}
	if !utils.IsValidPassword(req.NewPassword) {
		middleware.JSONError(c, "新密码需8-72位，包含大小写字母和数字喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)

	success, errMsg := service.ChangePassword(user.Uid, req.OldPassword, req.NewPassword)
	if !success {
		middleware.JSONError(c, errMsg, nil)
		return
	}

	middleware.JSONSuccess(c, "密码修改成功喵", nil)
}

func UpdateProfile(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	user := middleware.GetCurrentUserFromContext(c)

	allowedKeys := map[string]int{
		models.MetaGender:   0,
		models.MetaBirthday: 0,
		models.MetaBio:      200,
	}

	pendingFields := map[string]string{}

	for key, maxLen := range allowedKeys {
		if value, ok := req[key]; ok {
			strValue, _ := value.(string)
			if maxLen > 0 && len(strValue) > maxLen {
				middleware.JSONError(c, key+"长度不能超过"+strconv.Itoa(maxLen)+"个字符喵", nil)
				return
			}
			if key == models.MetaGender {
				gender, _ := strconv.Atoi(strValue)
				if gender < 0 || gender > 3 {
					middleware.JSONError(c, "性别值无效喵", nil)
					return
				}
			}
			pendingFields[key] = strValue
		}
	}

	if nameValue, ok := req["name"].(string); ok {
		if len(nameValue) > 50 {
			middleware.JSONError(c, "昵称长度不能超过50个字符喵", nil)
			return
		}
		pendingFields["name"] = nameValue
	}

	if len(pendingFields) > 0 {
		service.CreateUserMetaReview(user.Uid, pendingFields)
		middleware.JSONSuccess(c, "资料已提交审核喵", nil)
		return
	}

	middleware.JSONError(c, "没有需要更新的内容喵", nil)
}

func GetAllUsers(c *gin.Context) {
	users := service.GetUsersList()
	middleware.JSONSuccess(c, "ok", users)
}

func verifyCode(email, code string) (bool, string) {
	codeKey := db.FormatKey("verify:code", email)
	storedHash, err := db.Get(codeKey)

	if err != nil || storedHash == "" {
		return false, "验证码不存在或已过期喵喵"
	}
	if !utils.VerifyPassword(storedHash, code) {
		return false, "验证码错误喵喵"
	}
	db.Del(codeKey)
	return true, ""
}
