import type { UserInfo } from '~/types/domain'
import { authApi } from '~/api'
import { getRoleName } from '~/utils/mapping'

let pendingFetch: Promise<UserInfo | null> | null = null

export function useAuth() {
	const user = useState<UserInfo | null>('auth:user', () => null)
	const checked = useState<boolean>('auth:checked', () => false)

	const canAccessAdmin = computed(() => user.value?.canAccessAdmin === true)

	const displayName = computed(() => user.value?.name || '管理员')

	const identity = computed(() => {
		const roles = user.value?.roles || []
		if (roles.includes('superadmin')) return '超级管理员'
		if (roles.includes('admin')) return '管理员'
		return getRoleName(roles[0] || 'user')
	})

	async function fetchMe(force = false): Promise<UserInfo | null> {
		if (checked.value && !force) return user.value
		if (pendingFetch && !force) return pendingFetch

		pendingFetch = (async () => {
			try {
				const me = await authApi.me()
				user.value = me
				return me
			} catch {
				user.value = null
				return null
			} finally {
				checked.value = true
				pendingFetch = null
			}
		})()

		return pendingFetch
	}

	async function login(email: string, passwd: string): Promise<UserInfo | null> {
		await authApi.login(email, passwd)
		checked.value = false
		return fetchMe(true)
	}

	async function logout() {
		try {
			await authApi.logout()
		} catch {
			// 服务端会话可能已失效，本地状态仍需清理
		}
		user.value = null
		checked.value = true
	}

	return {
		user,
		checked,
		canAccessAdmin,
		displayName,
		identity,
		fetchMe,
		login,
		logout,
	}
}
