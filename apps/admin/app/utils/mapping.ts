export type MaoBadgeTone = 'neutral' | 'accent' | 'danger' | 'sage' | 'blue' | 'solid'

const STATUS_TONES: Record<string, MaoBadgeTone> = {
	publish: 'sage',
	draft: 'accent',
	private: 'blue',
	trash: 'danger',
	rejected: 'danger',
	approved: 'sage',
	pending: 'accent',
	spam: 'danger',
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

const ROLE_TONES: Record<string, MaoBadgeTone> = {
	superadmin: 'danger',
	admin: 'accent',
	editor: 'blue',
	user: 'sage',
	visitor: 'neutral',
}

const ROLE_NAMES: Record<string, string> = {
	superadmin: '超级管理员',
	admin: '管理员',
	editor: '编辑',
	user: '用户',
	visitor: '访客',
}

export const getStatusTone = (status: string): MaoBadgeTone => STATUS_TONES[status] ?? 'neutral'

export const getStatusName = (status: string): string => STATUS_NAMES[status] ?? status

export const getTypeName = (type: string): string => TYPE_NAMES[type] ?? type

export const getRoleTone = (role: string): MaoBadgeTone => ROLE_TONES[role] ?? 'neutral'

export const getRoleName = (role: string): string => ROLE_NAMES[role] ?? role
