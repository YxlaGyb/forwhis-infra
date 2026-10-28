package handler

import (
	"strconv"
	"strings"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/internal/request"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/db"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

const codeExpirationTime = 300

func Login(c fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	clientIP := middleware.ClientIP(c)
	if !utils.IsValidEmail(req.Email) || !utils.IsValidPassword(req.Passwd) {
		return middleware.JSONError(c, "邮箱或密码格式不正确喵", nil)
	}

	cdKey := db.FormatKey("login:cd", req.Email)
	if exists, _ := db.Exists(cdKey); exists > 0 {
		return middleware.JSONError(c, "操作过于频繁, 请稍后再试喵", nil)
	}
	db.SetEX(cdKey, "1", 60)

	result, errCode := service.LoginUser(req.Email, req.Passwd, clientIP)
	if errCode != 0 {
		db.Del(cdKey)
		return middleware.JSONError(c, "邮箱或密码错误喵", nil)
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	return middleware.JSONSuccess(c, "登录成功喵", fiber.Map{
		"uid": result["uid"], "name": result["name"], "email": result["email"],
	})

}

func Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "forestwhisper",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})
	return middleware.JSONSuccess(c, "退出成功喵", nil)
}

func Register(c fiber.Ctx) error {
	if !service.GetOptionBool("feature_enableRegistration", true) {
		return middleware.JSONError(c, "当前系统已关闭用户注册功能喵", nil)
	}

	var req request.RegisterRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if !utils.IsValidEmail(req.Email) || !utils.IsValidName(req.Username) || !utils.IsValidPassword(req.Passwd) {
		return middleware.JSONError(c, "邮箱、用户名或密码格式不正确喵（密码需8-16位, 包含字母和数字）", nil)
	}
	if len(req.Code) != 6 {
		return middleware.JSONError(c, "验证码应为6位字母+数字喵喵", nil)
	}
	if service.GetUserByEmail(req.Email) != nil {
		return middleware.JSONError(c, "该邮箱已被注册喵喵", nil)
	}

	valid, errMsg := verifyCode(req.Email, req.Code)
	if !valid {
		return middleware.JSONError(c, errMsg, nil)
	}

	clientIP := middleware.ClientIP(c)
	hashedPwd, err := utils.HashPassword(req.Passwd)
	if err != nil {
		return middleware.JSONSystemError(c, "密码加密失败喵喵", nil)
	}
	user, err := service.CreateUser(
		"",
		req.Email,
		hashedPwd,
		clientIP, "user",
		true)

	if err != nil {
		return middleware.JSONSystemError(c, "注册失败喵喵", nil)
	}

	if req.Username != "" {
		service.CreateUserMetaReview(user.Uid, map[string]string{"name": req.Username})
	}

	result, errCode := service.LoginUser(req.Email, req.Passwd, clientIP)
	if errCode != 0 {
		return middleware.JSONError(c, "注册后登录失败喵", nil)
	}

	token := result["token"].(string)
	utils.SetCookieToken(c, token)
	return middleware.JSONSuccess(c, "注册成功喵", nil)
}

func SendEmailCodeRegister(c fiber.Ctx) error {
	if !service.GetOptionBool("feature_enableRegistration", true) {
		return middleware.JSONError(c, "当前系统已关闭用户注册功能喵", nil)
	}

	var req request.SendEmailCodeRequest

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	if !utils.IsValidEmail(req.Email) {
		return middleware.JSONError(c, "邮箱格式错误喵", nil)
	}
	if service.GetUserByEmail(req.Email) != nil {
		return middleware.JSONError(c, "您的邮箱已经被使用了喵, 请换一个试试喵", nil)
	}

	code := utils.SecureRandomString(6)
	hashedCode, err := utils.HashPassword(code)
	if err != nil {
		return middleware.JSONSystemError(c, "验证码生成失败喵喵", nil)
	}
	htmlContent := `
		<div style="font-family: sans-serif; padding: 20px; border: 1px solid #eee; border-radius: 8px;">
			<h2 style="color: #333;">✨ 注册验证码待查收</h2>
			<p>您的验证码是 <b style="font-size: 24px; color: #007bff;">` + code + `</b></p>
			<p style="color: #666;">请在 5 分钟内完成验证, 请勿泄露喵~</p>
			<hr style="border: none; border-top: 1px solid #eee;">
			<footer style="font-size: 12px; color: #999;">来自：数数洞洞平台</footer>
		</div>
	`

	if !utils.SendMail("注册验证码", []string{req.Email}, htmlContent) {
		return middleware.JSONSystemError(c, "邮件服务连接超时, 请稍后再试喵喵", nil)
	}

	codeKey := db.FormatKey("verify:code", req.Email)
	db.SetEX(codeKey, hashedCode, codeExpirationTime)
	return middleware.JSONSuccess(c, "验证码已发送喵", fiber.Map{"codeSalt": hashedCode})
}

func GetUserByUsername(c fiber.Ctx) error {
	name := c.Query("name")
	exactly := c.Query("exactly")

	if strings.TrimSpace(name) == "" {
		return middleware.JSONError(c, "用户名不能为空喵", nil)
	}

	users := service.GetUsersByUsername(name, exactly)
	return middleware.JSONSuccess(c, "ok", users)
}

func DeleteUserByUid(c fiber.Ctx) error {
	var req request.DeleteUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	if req.Uid == 0 {
		return middleware.JSONError(c, "用户ID不能为空喵", nil)
	}

	currentUser := middleware.GetCurrentUserFromContext(c)
	if currentUser.Uid == req.Uid {
		return middleware.JSONError(c, "不能删除自己喵", nil)
	}

	userWithRoles, ok := service.GetUserWithRoles(req.Uid)
	if !ok {
		return middleware.JSONError(c, "用户不存在喵喵", nil)
	}

	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)
	isTargetSuperadmin := rbacSvc.IsSuperadmin(userWithRoles.User)
	if isTargetSuperadmin && !isCurrentSuperadmin {
		return middleware.JSONError(c, "无权删除超级管理员喵", nil)
	}

	success := service.DeleteUser(currentUser.Uid, req.Uid)
	if !success {
		return middleware.JSONError(c, "用户删除失败喵喵", nil)
	}
	return middleware.JSONSuccess(c, "用户删除成功喵", nil)
}

func ChangePassword(c fiber.Ctx) error {
	var req request.ChangePasswordRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return middleware.JSONError(c, "旧密码和新密码不能为空喵", nil)
	}
	if !utils.IsValidPassword(req.NewPassword) {
		return middleware.JSONError(c, "新密码需8-16位, 包含字母和数字喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)

	success, errMsg := service.ChangePassword(user.Uid, req.OldPassword, req.NewPassword)
	if !success {
		return middleware.JSONError(c, errMsg, nil)
	}
	return middleware.JSONSuccess(c, "密码修改成功喵", nil)
}

func UpdateProfile(c fiber.Ctx) error {
	var req map[string]interface{}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)

	allowedKeys := map[string]int{
		models.ProfileGender:   0,
		models.ProfileBirthday: 10,
		models.ProfileBio:      200,
	}

	pendingFields := map[string]string{}

	for key, maxLen := range allowedKeys {
		if value, ok := req[key]; ok {
			strValue, _ := value.(string)
			if maxLen > 0 && len(strValue) > maxLen {
				return middleware.JSONError(c, key+"长度不能超过"+strconv.Itoa(maxLen)+"个字符喵", nil)
			}
			if key == models.ProfileGender {
				gender, _ := strconv.Atoi(strValue)
				if gender < 0 || gender > 3 {
					return middleware.JSONError(c, "性别值无效喵", nil)
				}
			}
			pendingFields[key] = strValue
		}
	}

	if nameValue, ok := req["name"].(string); ok {
		if len(nameValue) > 50 {
			return middleware.JSONError(c, "昵称长度不能超过50个字符喵", nil)
		}
		pendingFields["name"] = nameValue
	}

	if len(pendingFields) > 0 {
		service.CreateUserMetaReview(user.Uid, pendingFields)
		return middleware.JSONSuccess(c, "资料已提交审核喵", nil)
	}
	return middleware.JSONError(c, "没有需要更新的内容喵", nil)
}

func GetAllUsers(c fiber.Ctx) error {
	users := service.GetUsersList()
	return middleware.JSONSuccess(c, "ok", users)
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
