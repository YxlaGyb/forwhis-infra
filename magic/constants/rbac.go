package constants

type Role string

const (
	RoleSuperadmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleEditor     Role = "editor"
	RoleUser       Role = "user"
	RoleVisitor    Role = "visitor"
)

type Permission string

const (
	PermUserRead        Permission = "user:read"
	PermUserCreate      Permission = "user:create"
	PermUserEdit        Permission = "user:edit"
	PermUserDelete      Permission = "user:delete"
	PermPostRead        Permission = "post:read"
	PermPostCreate      Permission = "post:create"
	PermPostEdit        Permission = "post:edit"
	PermPostDelete      Permission = "post:delete"
	PermPostPublish     Permission = "post:publish"
	PermCommentRead     Permission = "comment:read"
	PermCommentCreate   Permission = "comment:create"
	PermCommentDelete   Permission = "comment:delete"
	PermCommentModerate Permission = "comment:moderate"
	PermSystemConfig    Permission = "system:config"
	PermSystemLogs      Permission = "system:logs"
	PermSystemBackup    Permission = "system:backup"
	PermSystemPlugin    Permission = "system:plugin"
	PermSettingsEdit    Permission = "settings:edit"
)

type PermissionCategory string

const (
	CatUserMgmt    PermissionCategory = "用户管理"
	CatContentMgmt PermissionCategory = "内容管理"
	CatSystemMgmt  PermissionCategory = "系统管理"
)

type PermDef struct {
	Name     Permission
	Desc     string
	Category PermissionCategory
}

var AllPermissions = []PermDef{
	{PermUserRead, "查看用户信息", CatUserMgmt},
	{PermUserCreate, "创建新用户", CatUserMgmt},
	{PermUserEdit, "编辑用户信息", CatUserMgmt},
	{PermUserDelete, "删除用户", CatUserMgmt},
	{PermPostRead, "查看文章", CatContentMgmt},
	{PermPostCreate, "创建文章", CatContentMgmt},
	{PermPostEdit, "编辑文章", CatContentMgmt},
	{PermPostDelete, "删除文章", CatContentMgmt},
	{PermPostPublish, "发布文章", CatContentMgmt},
	{PermCommentRead, "查看评论", CatContentMgmt},
	{PermCommentCreate, "创建评论", CatContentMgmt},
	{PermCommentDelete, "删除评论", CatContentMgmt},
	{PermCommentModerate, "审核评论", CatContentMgmt},
	{PermSystemConfig, "系统配置", CatSystemMgmt},
	{PermSystemLogs, "查看日志", CatSystemMgmt},
	{PermSystemBackup, "数据备份", CatSystemMgmt},
	{PermSystemPlugin, "插件管理", CatSystemMgmt},
	{PermSettingsEdit, "编辑系统设置", CatSystemMgmt},
}

var RolePermissions = map[Role][]Permission{
	RoleSuperadmin: {PermUserRead, PermUserCreate, PermUserEdit, PermUserDelete, PermPostRead, PermPostCreate, PermPostEdit, PermPostDelete, PermPostPublish, PermCommentRead, PermCommentCreate, PermCommentDelete, PermCommentModerate, PermSystemConfig, PermSystemLogs, PermSystemBackup, PermSystemPlugin, PermSettingsEdit},
	RoleAdmin:      {PermUserRead, PermUserCreate, PermUserEdit, PermUserDelete, PermPostRead, PermPostCreate, PermPostEdit, PermPostDelete, PermPostPublish, PermCommentRead, PermCommentCreate, PermCommentDelete, PermCommentModerate, PermSystemLogs, PermSettingsEdit},
	RoleEditor:     {PermPostRead, PermPostCreate, PermPostEdit, PermPostPublish, PermCommentRead, PermCommentCreate},
	RoleUser:       {PermPostRead, PermCommentRead, PermCommentCreate},
	RoleVisitor:    {PermPostRead},
}

var RoleDescriptions = map[Role]string{
	RoleSuperadmin: "超级管理员，拥有系统所有权限",
	RoleAdmin:      "管理员，负责系统日常管理",
	RoleEditor:     "编辑者，负责内容创作和编辑",
	RoleUser:       "普通用户，可阅读内容和发表评论",
	RoleVisitor:    "访客，仅可浏览公开内容",
}
