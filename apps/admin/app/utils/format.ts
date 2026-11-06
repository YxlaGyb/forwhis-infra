import { TIMESTAMP_LIMITS } from './constants'

const toDate = (timestamp: number): Date | null => {
	if (!timestamp || typeof timestamp !== 'number') return null
	if (!Number.isFinite(timestamp)) return null
	if (timestamp < TIMESTAMP_LIMITS.minValid || timestamp > TIMESTAMP_LIMITS.maxValid) return null
	const date = new Date(timestamp * 1000)
	return Number.isNaN(date.getTime()) ? null : date
}

export const formatTimestamp = (timestamp: number): string => {
	const date = toDate(timestamp)
	if (!date) return '-'
	return date.toLocaleString('zh-CN', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit',
	})
}

export const formatDate = (timestamp: number): string => {
	const date = toDate(timestamp)
	if (!date) return '-'
	return date.toLocaleDateString('zh-CN', {
		year: 'numeric',
		month: '2-digit',
		day: '2-digit',
	})
}
