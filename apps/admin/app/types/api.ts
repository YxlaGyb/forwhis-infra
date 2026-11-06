export interface ApiResponse<T = unknown> {
	code: number
	msg?: string
	data?: T
}

export interface Paginated<T> {
	total: number
	page: number
	limit: number
	items: T[]
}
