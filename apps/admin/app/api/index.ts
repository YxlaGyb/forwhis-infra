import type {
	AdminUser,
	AllSettings,
	Category,
	Comment,
	Content,
	DashboardData,
	LoginResult,
	TrendPoint,
	UserInfo,
} from '~/types/domain'
import type { Paginated } from '~/types/api'
import { request } from './client'

export const authApi = {
	me: () => request<UserInfo>('/api/v1/user/me'),
	login: (email: string, passwd: string) =>
		request<LoginResult>('/api/v1/user/login', { method: 'POST', body: { email, passwd } }),
	logout: () => request<null>('/api/v1/user/logout', { method: 'POST' }),
}

export const statsApi = {
	dashboard: () => request<DashboardData>('/api/v1/admin/stats/dashboard'),
	trend: (days: number) => request<TrendPoint[]>('/api/v1/admin/stats/trend', { query: { days } }),
}

export interface CreateUserPayload {
	name: string
	email: string
	password: string
	roles: string[]
}

export interface UpdateUserPayload {
	name?: string
	email?: string
	roles?: string[]
}

export const usersApi = {
	list: () => request<AdminUser[]>('/api/v1/admin/users'),
	create: (data: CreateUserPayload) =>
		request<{ uid: number }>('/api/v1/admin/user/create', { method: 'POST', body: data }),
	update: (uid: number, data: UpdateUserPayload) =>
		request<null>(`/api/v1/admin/user/${uid}/update`, { method: 'POST', body: data }),
	remove: (uid: number) =>
		request<null>('/api/v1/user/users/del', { method: 'POST', body: { uid } }),
}

export interface ContentListParams {
	page?: number
	limit?: number
	status?: string
	type?: string
}

export interface ContentPayload {
	title?: string
	text?: string
	status?: string
	type?: string
	categoryMid?: number
}

export const contentsApi = {
	list: (params?: ContentListParams) =>
		request<Paginated<Content> | Content[]>('/api/v1/content', { query: { ...params } }),
	get: (cid: number) => request<Content>(`/api/v1/content/${cid}`),
	create: (data: ContentPayload) =>
		request<{ cid: number }>('/api/v1/content', { method: 'POST', body: data }),
	update: (cid: number, data: ContentPayload) =>
		request<null>(`/api/v1/content/${cid}`, { method: 'POST', body: data }),
	remove: (cid: number) => request<null>(`/api/v1/content/${cid}/del`, { method: 'POST' }),
}

export interface CommentListParams {
	page?: number
	limit?: number
	status?: string
	cid?: number
}

export const commentsApi = {
	list: (params?: CommentListParams) =>
		request<Paginated<Comment> | Comment[]>('/api/v1/comment', { query: { ...params } }),
	update: (coid: number, data: { text?: string; status?: string }) =>
		request<Comment>(`/api/v1/comment/${coid}`, { method: 'POST', body: data }),
	approve: (coid: number) => request<null>(`/api/v1/comment/${coid}/approve`, { method: 'POST' }),
	remove: (coid: number) => request<null>(`/api/v1/comment/${coid}/del`, { method: 'POST' }),
}

export interface MetaPayload {
	name?: string
	slug?: string
	description?: string
	parent?: number
	type?: string
}

export const metasApi = {
	list: (type: string) => request<Category[]>('/api/v1/meta', { query: { type } }),
	create: (data: MetaPayload) =>
		request<{ mid: number }>('/api/v1/meta', { method: 'POST', body: data }),
	update: (mid: number, data: MetaPayload) =>
		request<Category>(`/api/v1/meta/${mid}`, { method: 'POST', body: data }),
	remove: (mid: number) => request<null>(`/api/v1/meta/${mid}/del`, { method: 'POST' }),
}

export const settingsApi = {
	get: () => request<Partial<AllSettings>>('/api/v1/admin/settings'),
	update: (data: AllSettings) =>
		request<AllSettings>('/api/v1/admin/settings', { method: 'POST', body: data }),
}

export function listItems<T>(payload: T[] | Paginated<T> | undefined | null): T[] {
	if (!payload) return []
	if (Array.isArray(payload)) return payload
	return payload.items || []
}
