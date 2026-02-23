package controller

import (
	"strconv"
	"time"

	rbacSvc "forestwhisperGo/magic/service/rbac"

	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/models"
	"forestwhisperGo/magic/models/types"
	"forestwhisperGo/magic/service"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

	"github.com/gin-gonic/gin"
)

func GetPendingReviews(c *gin.Context) {
	var reviews []models.Review
	db.GetDB().Where("status = ?", models.ReviewStatusPending).Find(&reviews)

	data := []map[string]interface{}{}
	for _, r := range reviews {
		data = append(data, map[string]interface{}{
			"review_id":   r.ReviewId,
			"target_type": r.TargetType,
			"target_id":   r.TargetId,
			"created_at":  r.CreatedAt,
		})
	}
	middleware.JSONSuccess(c, "获取成功喵", data)
}

func ApproveReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	var req types.ApproveReviewRequest
	c.ShouldBindJSON(&req)

	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		middleware.JSONError(c, "审核记录不存在喵", nil)
		return
	}

	if review.TargetType == models.TargetTypeUserProfile {
		var meta models.Usermeta
		if db.GetDB().First(&meta, review.TargetId).Error == nil && req.NewValue != "" {
			meta.MetaValue = req.NewValue
			meta.UpdatedAt = int(time.Now().Unix())
			db.GetDB().Save(&meta)
		}
	}

	review.Status = models.ReviewStatusApproved
	review.UpdatedAt = int(time.Now().Unix())
	user := middleware.GetCurrentUser(c)
	if user != nil {
		review.ReviewerId = user.Uid
	}
	db.GetDB().Save(&review)

	if user != nil {
		service.CreateUserinfo(user.Uid, "approveReview", "通过了审核: "+strconv.Itoa(reviewId)+"喵")
	}

	middleware.JSONSuccess(c, "审核已通过喵", nil)
}

func RejectReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		middleware.JSONError(c, "审核ID无效喵", nil)
		return
	}

	var review models.Review
	if db.GetDB().First(&review, reviewId).Error != nil {
		middleware.JSONError(c, "审核记录不存在喵", nil)
		return
	}

	var req types.RejectReviewRequest
	c.ShouldBindJSON(&req)

	review.Status = models.ReviewStatusRejected
	review.UpdatedAt = int(time.Now().Unix())
	review.Reason = req.Reason
	user := middleware.GetCurrentUser(c)
	if user != nil {
		review.ReviewerId = user.Uid
	}
	db.GetDB().Save(&review)

	if user != nil {
		service.CreateUserinfo(user.Uid, "rejectReview", "拒绝了审核: "+strconv.Itoa(reviewId)+"喵")
	}

	middleware.JSONSuccess(c, "审核已拒绝喵", nil)
}

func AdminUpdateUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		middleware.JSONError(c, "用户ID无效喵", nil)
		return
	}

	var req types.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请求数据不能为空喵", nil)
		return
	}

	req.Uid = uid

	currentUser := middleware.GetCurrentUserFromContext(c)
	targetUserWithRoles, ok := service.GetUserWithRoles(uid)
	if !ok {
		middleware.JSONError(c, "用户不存在喵", nil)
		return
	}

	isCurrentSuperadmin := rbacSvc.IsSuperadmin(currentUser)
	isTargetSuperadmin := rbacSvc.IsSuperadmin(targetUserWithRoles.User)
	isCurrentAdmin := rbacSvc.HasAnyAdminRole(currentUser)

	if req.Name != "" && !utils.IsValidName(req.Name) {
		middleware.JSONError(c, "昵称格式不正确喵", nil)
		return
	}

	if req.Email != "" && !utils.IsValidEmail(req.Email) {
		middleware.JSONError(c, "邮箱格式不正确喵", nil)
		return
	}

	if len(req.Roles) > 0 {
		if !isCurrentSuperadmin && !isCurrentAdmin {
			middleware.JSONError(c, "您没有权限修改用户角色喵", nil)
			return
		}

		hasSuperadminRole := false
		for _, r := range req.Roles {
			if r == "superadmin" {
				hasSuperadminRole = true
				break
			}
		}

		if hasSuperadminRole && !isCurrentSuperadmin {
			middleware.JSONError(c, "只有超级管理员可以授予超级管理员权限喵", nil)
			return
		}

		if currentUser.Uid == uid && isCurrentSuperadmin && !hasSuperadminRole {
			middleware.JSONError(c, "超级管理员不能移除自己的超级管理员权限喵", nil)
			return
		}

		if !isCurrentSuperadmin && isTargetSuperadmin {
			middleware.JSONError(c, "无权修改超级管理员喵", nil)
			return
		}

		if !isCurrentSuperadmin && isCurrentAdmin {
			if rbacSvc.IsSuperadmin(targetUserWithRoles.User) {
				middleware.JSONError(c, "无权修改超级管理员喵", nil)
				return
			}
			isTargetAdmin := false
			for _, r := range targetUserWithRoles.Roles {
				if r.Role.Name == "admin" {
					isTargetAdmin = true
					break
				}
			}
			if isTargetAdmin && currentUser.Uid != uid {
				middleware.JSONError(c, "管理员不能修改其他管理员的权限喵", nil)
				return
			}
		}
	}

	success := service.AdminUpdateUser(uid, req.Name, req.Email, req.Roles, req.IsActive)
	if !success {
		middleware.JSONError(c, "更新用户失败喵", nil)
		return
	}

	rbacSvc.InvalidateUserCache(uid)
	service.CreateUserinfo(currentUser.Uid, "updateUser", "更新了用户: "+strconv.Itoa(uid)+"喵")
	middleware.JSONSuccess(c, "用户更新成功喵", nil)
}
