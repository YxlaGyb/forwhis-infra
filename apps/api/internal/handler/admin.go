package handler

import (
	"strconv"

	"forestwhisperGo/apps/api/internal/middleware"
	"forestwhisperGo/apps/api/internal/request"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"
	"forestwhisperGo/apps/api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func AdminCreateUser(c fiber.Ctx) error {
	var req request.AdminCreateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	if !utils.IsValidName(req.Name) {
		return middleware.JSONError(c, "昵称格式不正确喵", nil)
	}

	if !utils.IsValidEmail(req.Email) {
		return middleware.JSONError(c, "邮箱格式不正确喵", nil)
	}

	if !utils.IsValidPassword(req.Password) {
		return middleware.JSONError(c, "密码需8-16位, 包含字母和数字喵", nil)
	}

	if service.GetUserByEmail(req.Email) != nil {
		return middleware.JSONError(c, "该邮箱已被使用喵", nil)
	}

	for _, r := range req.Roles {
		if !constants.IsValidRole(r) {
			return middleware.JSONError(c, "角色无效喵", nil)
		}
	}

	currentUser := middleware.GetCurrentUserFromContext(c)
	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)

	if len(req.Roles) > 0 {
		hasSuperadminRole := false
		for _, r := range req.Roles {
			if r == "superadmin" {
				hasSuperadminRole = true
				break
			}
		}
		if hasSuperadminRole && !isCurrentSuperadmin {
			return middleware.JSONError(c, "只有超级管理员可以授予超级管理员权限喵", nil)
		}
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return middleware.JSONSystemError(c, "密码加密失败喵", nil)
	}

	role := "user"
	if len(req.Roles) > 0 {
		role = req.Roles[0]
	}

	clientIP := middleware.ClientIP(c)
	user, err := service.CreateUser(req.Name, req.Email, hashedPwd, clientIP, role, false)
	if err != nil {
		return middleware.JSONSystemError(c, "创建用户失败喵", nil)
	}

	if len(req.Roles) > 1 {
		for _, r := range req.Roles[1:] {
			rbacSvc.AssignRoleToUser(user.Uid, r, currentUser.Uid)
		}
	}

	rbacSvc.InvalidateUserCache(user.Uid)
	service.CreateUserinfo(currentUser.Uid, "createUser", "创建了用户: "+req.Email+"喵")
	return middleware.JSONSuccess(c, "用户创建成功喵", fiber.Map{"uid": user.Uid})
}

func AdminUpdateUser(c fiber.Ctx) error {
	uid, err := strconv.Atoi(c.Params("uid"))
	if err != nil {
		return middleware.JSONError(c, "用户ID无效喵", nil)
	}

	var req request.AdminUpdateUserRequest
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}

	req.Uid = uid

	currentUser := middleware.GetCurrentUserFromContext(c)
	targetUserWithRoles, ok := service.GetUserWithRoles(uid)
	if !ok {
		return middleware.JSONError(c, "用户不存在喵", nil)
	}

	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)
	isTargetSuperadmin := rbacSvc.IsSuperadmin(targetUserWithRoles.User)
	isCurrentAdmin := rbacSvc.HasAnyAdminRole(currentUser)

	if req.Name != "" && !utils.IsValidName(req.Name) {
		return middleware.JSONError(c, "昵称格式不正确喵", nil)
	}

	if req.Email != "" && !utils.IsValidEmail(req.Email) {
		return middleware.JSONError(c, "邮箱格式不正确喵", nil)
	}

	for _, r := range req.Roles {
		if !constants.IsValidRole(r) {
			return middleware.JSONError(c, "角色无效喵", nil)
		}
	}

	if len(req.Roles) > 0 {
		if !isCurrentSuperadmin && !isCurrentAdmin {
			return middleware.JSONError(c, "您没有权限修改用户角色喵", nil)
		}

		hasSuperadminRole := false
		for _, r := range req.Roles {
			if r == "superadmin" {
				hasSuperadminRole = true
				break
			}
		}

		if hasSuperadminRole && !isCurrentSuperadmin {
			return middleware.JSONError(c, "只有超级管理员可以授予超级管理员权限喵", nil)
		}

		if currentUser.Uid == uid && isCurrentSuperadmin && !hasSuperadminRole {
			return middleware.JSONError(c, "超级管理员不能移除自己的超级管理员权限喵", nil)
		}

		if !isCurrentSuperadmin && isTargetSuperadmin {
			return middleware.JSONError(c, "无权修改超级管理员喵", nil)
		}

		if !isCurrentSuperadmin && isCurrentAdmin {
			isTargetAdmin := false
			for _, r := range targetUserWithRoles.Roles {
				if r == "admin" {
					isTargetAdmin = true
					break
				}
			}
			if isTargetAdmin && currentUser.Uid != uid {
				return middleware.JSONError(c, "管理员不能修改其他管理员的权限喵", nil)
			}
		}
	}

	success := service.AdminUpdateUser(uid, req.Name, req.Email, req.Roles, req.IsActive)
	if !success {
		return middleware.JSONError(c, "更新用户失败喵", nil)
	}

	rbacSvc.InvalidateUserCache(uid)
	service.CreateUserinfo(currentUser.Uid, "updateUser", "更新了用户: "+strconv.Itoa(uid)+"喵")
	return middleware.JSONSuccess(c, "用户更新成功喵", nil)
}
