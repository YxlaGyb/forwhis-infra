import type { ApiResponse } from '~/types/api'

export interface ApiRequestOptions {
	method?: 'GET' | 'POST'
	body?: unknown
	query?: Record<string, unknown>
}

export type ApiClient = <T>(url: string, options?: ApiRequestOptions) => Promise<T>

let client: ApiClient | null = null

export function setApiClient(instance: ApiClient) {
	client = instance
}

export class ApiError extends Error {
	readonly code: number

	constructor(message: string, code = 0) {
		super(message)
		this.name = 'ApiError'
		this.code = code
	}
}

export function unwrap<T>(response: ApiResponse<T> | undefined): T {
	if (response?.code === 200) return response.data as T
	throw new ApiError(response?.msg || '请求失败', response?.code ?? 0)
}

export async function request<T>(url: string, options?: ApiRequestOptions): Promise<T> {
	if (!client) throw new ApiError('API 客户端尚未初始化')
	return unwrap(await client<ApiResponse<T>>(url, options))
}

const AUTH_FAILURE_PATTERN = /未登录|登录已过期|登录失效|token/i

export function isAuthFailure(payload: unknown): boolean {
	if (!payload || typeof payload !== 'object') return false
	const response = payload as ApiResponse<unknown>
	return response.code !== 200 && AUTH_FAILURE_PATTERN.test(response.msg || '')
}

export function errorMessage(error: unknown, fallback = '请求失败'): string {
	if (error instanceof ApiError) return error.message
	const data = (error as { data?: { msg?: string } } | null)?.data
	if (data?.msg) return data.msg
	const message = (error as { message?: string } | null)?.message
	return message || fallback
}
