import type { ContentType } from '~/types/domain'

export const TIMESTAMP_LIMITS = {
	minValid: 0,
	maxValid: 2147483647,
} as const

export const CONTENT_LIMITS: Record<ContentType, number> = {
	post: 1000,
	article: 30000,
	page: 30000,
}

export const EDITOR_TITLE_MAX_LENGTH = 30

export const USER_LIMITS = {
	passwordMinLength: 8,
	passwordMaxLength: 16,
	passwordPattern: /^(?=.*[a-zA-Z])(?=.*\d)/,
	nameMaxLength: 50,
	emailMaxLength: 100,
} as const

export const DEFAULT_PAGE_SIZE = 10
export const PAGE_SIZE_OPTIONS = ['10', '20', '50', '100']

export const ROLE_OPTIONS = [
	{ value: 'user', label: '普通用户' },
	{ value: 'editor', label: '编辑' },
	{ value: 'admin', label: '管理员' },
	{ value: 'superadmin', label: '超级管理员' },
]

export const SETTINGS_SECTIONS = [
	{ key: 'site', title: '站点信息' },
	{ key: 'feature', title: '功能开关' },
	{ key: 'seo', title: 'SEO设置' },
] as const
