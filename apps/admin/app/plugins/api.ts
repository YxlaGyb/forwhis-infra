import type { UserInfo } from '~/types/domain'
import type { ApiClient } from '~/api/client'
import { isAuthFailure, setApiClient } from '~/api/client'

export default defineNuxtPlugin(() => {
	const config = useRuntimeConfig()
	const router = useRouter()
	const user = useState<UserInfo | null>('auth:user', () => null)
	const checked = useState<boolean>('auth:checked', () => false)

	const handleAuthFailure = () => {
		user.value = null
		checked.value = true
		if (import.meta.server) return
		if (router.currentRoute.value.path !== '/login') {
			void router.push('/login')
		}
	}

	const api = $fetch.create({
		baseURL: import.meta.server
			? config.apiBaseSsr || config.public.apiBaseUrl
			: config.public.apiBaseUrl,
		credentials: 'include',
		onResponse({ response }) {
			if (isAuthFailure(response._data)) handleAuthFailure()
		},
		onResponseError({ response }) {
			if (response.status === 401 || response.status === 403) handleAuthFailure()
		},
	})

	setApiClient(api as unknown as ApiClient)

	return { provide: { api } }
})
