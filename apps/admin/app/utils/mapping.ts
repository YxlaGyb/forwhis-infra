const STATUS_COLORS: Record<string, string> = {
	publish: 'green',
	draft: 'orange',
	private: 'blue',
	trash: 'red',
	rejected: 'red',
	approved: 'green',
	pending: 'orange',
	spam: 'red',
}

const STATUS_NAMES: Record<string, string> = {
	publish: '已发布',
	draft: '草稿',
	private: '私密',
	trash: '回收站',
	rejected: '已拒绝',
	approved: '已批准',
	pending: '待审核',
	spam: '垃圾',
}

const TYPE_NAMES: Record<string, string> = {
	post: '帖子',
	article: '文章',
	page: '页面',
}

const ROLE_COLORS: Record<string, string> = {
	superadmin: 'red',
	admin: 'orange',
	editor: 'blue',
	user: 'green',
	visitor: 'default',
}

const ROLE_NAMES: Record<string, string> = {
	superadmin: '超级管理员',
	admin: '管理员',
	editor: '编辑',
	user: '用户',
	visitor: '访客',
}

export const getStatusColor = (status: string): string => STATUS_COLORS[status] || 'default'

export const getStatusName = (status: string): string => STATUS_NAMES[status] || status

export const getTypeName = (type: string): string => TYPE_NAMES[type] || type

export const getRoleColor = (role: string): string => ROLE_COLORS[role] || 'default'

export const getRoleName = (role: string): string => ROLE_NAMES[role] || role
