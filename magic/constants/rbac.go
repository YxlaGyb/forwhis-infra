package constants

const (
	RoleSuperadmin = "superadmin"
	RoleAdmin      = "admin"
	RoleEditor     = "editor"
	RoleUser       = "user"
	RoleVisitor    = "visitor"
)

const (
	CatUserMgmt    = "用户管理"
	CatContentMgmt = "内容管理"
	CatSystemMgmt  = "系统管理"
)

const (
	UserRead        = "user:read"
	UserCreate      = "user:create"
	UserEdit        = "user:edit"
	UserDelete      = "user:delete"
	PostRead        = "post:read"
	PostCreate      = "post:create"
	PostEdit        = "post:edit"
	PostDelete      = "post:delete"
	PostPublish     = "post:publish"
	CommentRead     = "comment:read"
	CommentCreate   = "comment:create"
	CommentDelete   = "comment:delete"
	CommentModerate = "comment:moderate"
	SystemConfig    = "system:config"
	SystemLogs      = "system:logs"
	SystemBackup    = "system:backup"
	SystemPlugin    = "system:plugin"
)

type PermDef struct {
	Name     string
	Desc     string
	Category string
}

var AllPermissions = []PermDef{
	{UserRead, "查看用户信息", CatUserMgmt},
	{UserCreate, "创建新用户", CatUserMgmt},
	{UserEdit, "编辑用户信息", CatUserMgmt},
	{UserDelete, "删除用户", CatUserMgmt},
	{PostRead, "查看文章", CatContentMgmt},
	{PostCreate, "创建文章", CatContentMgmt},
	{PostEdit, "编辑文章", CatContentMgmt},
	{PostDelete, "删除文章", CatContentMgmt},
	{PostPublish, "发布文章", CatContentMgmt},
	{CommentRead, "查看评论", CatContentMgmt},
	{CommentCreate, "创建评论", CatContentMgmt},
	{CommentDelete, "删除评论", CatContentMgmt},
	{CommentModerate, "审核评论", CatContentMgmt},
	{SystemConfig, "系统配置", CatSystemMgmt},
	{SystemLogs, "查看日志", CatSystemMgmt},
	{SystemBackup, "数据备份", CatSystemMgmt},
	{SystemPlugin, "插件管理", CatSystemMgmt},
}

var RolePermissions = map[string][]string{
	RoleSuperadmin: {UserRead, UserCreate, UserEdit, UserDelete, PostRead, PostCreate, PostEdit, PostDelete, PostPublish, CommentRead, CommentCreate, CommentDelete, CommentModerate, SystemConfig, SystemLogs, SystemBackup, SystemPlugin},
	RoleAdmin:      {UserRead, UserCreate, UserEdit, UserDelete, PostRead, PostCreate, PostEdit, PostDelete, PostPublish, CommentRead, CommentCreate, CommentDelete, CommentModerate, SystemLogs},
	RoleEditor:     {PostRead, PostCreate, PostEdit, PostPublish, CommentRead, CommentCreate},
	RoleUser:       {PostRead, CommentRead, CommentCreate},
	RoleVisitor:    {PostRead},
}

var RoleDescriptions = map[string]string{
	RoleSuperadmin: "超级管理员，拥有系统所有权限",
	RoleAdmin:      "管理员，负责系统日常管理",
	RoleEditor:     "编辑者，负责内容创作和编辑",
	RoleUser:       "普通用户，可阅读内容和发表评论",
	RoleVisitor:    "访客，仅可浏览公开内容",
}
