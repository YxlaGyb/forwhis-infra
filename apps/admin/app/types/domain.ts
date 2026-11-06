import type { Component } from 'vue'

export type ContentType = 'post' | 'article' | 'page'
export type ContentStatus = 'publish' | 'draft' | 'private' | 'trash' | 'pending' | 'rejected'
export type CommentStatus = 'approved' | 'pending' | 'spam' | 'trash'

export interface UserInfo {
	uid: number
	name: string
	email: string
	roles: string[]
	permissions?: string[]
	canAccessAdmin?: boolean
	lastip?: string
	createdAt?: number
	lastLogin?: number
	isActive?: number
	bio?: string
	birthday?: string
	gender?: string
}

export interface LoginResult {
	uid: number
	name: string
	email: string
}

export interface AdminUser {
	uid: number
	name: string
	email: string
	roles: string[]
	createdAt: number
	lastLogin?: number
	lastip?: string
	isActive?: number
	canAccessAdmin?: boolean
}

export interface ContentAuthor {
	uid: number
	name: string
}

export interface Content {
	cid: number
	title: string
	slug?: string
	text: string
	status: string
	type: string
	author?: string | ContentAuthor
	author_id?: number
	comments_num?: number
	commentsNum?: number
	created: number
	modified?: number
	views?: number
	tags?: unknown[]
	category?: { mid: number; name: string }
	categoryMid?: number
	hasPassword?: boolean
}

export interface Comment {
	coid: number
	cid: number
	text: string
	author: string
	mail: string
	url?: string
	ip?: string
	status: string
	parent?: number
	created: number
}

export interface Category {
	mid: number
	name: string
	slug: string
	type?: string
	description?: string
	parent?: number
	count?: number
}

export interface UserStats {
	total: number
	todayNew: number
	weekNew: number
}

export interface ContentStats {
	posts: number
	articles: number
	pages: number
	drafts: number
	published: number
	private: number
	views: number
}

export interface CommentStats {
	total: number
	pending: number
	approved: number
	spam: number
}

export interface DashboardData {
	users?: UserStats
	contents?: ContentStats
	comments?: CommentStats
}

export interface TrendPoint {
	date: string
	user: number
	content: number
	comment: number
	views: number
}

export interface SiteSettings {
	siteName: string
	siteDescription: string
	siteLogo: string
	favicon: string
	icp: string
}

export interface FeatureSettings {
	enableRegistration: boolean
	enableComments: boolean
	enableEmailVerify: boolean
	enableCaptcha: boolean
}

export interface SeoSettings {
	siteKeywords: string
	siteAuthor: string
	robotsTxt: string
	enableSitemap: boolean
}

export interface AllSettings {
	site: SiteSettings
	feature: FeatureSettings
	seo: SeoSettings
}

export type SettingsSection = 'site' | 'feature' | 'seo'

export interface RoleDefinition {
	name: string
	description: string
	permissions: string[]
	userCount: number
}

export interface ActionItem {
	key: string
	text: string
	icon?: Component
	danger?: boolean
	disabled?: boolean
	hidden?: boolean
	onClick: () => void
}

export interface ResponsiveColumn {
	key: string
	title: string
	dataIndex?: string
	priority?: 'high' | 'medium' | 'low'
	ellipsis?: boolean
}
