export default defineNuxtRouteMiddleware(async (to) => {
	const { user, checked, canAccessAdmin, fetchMe } = useAuth()

	if (!checked.value) await fetchMe()

	if (!user.value) {
		return to.path === '/login' ? undefined : navigateTo('/login')
	}

	if (!canAccessAdmin.value) {
		return to.path === '/no-access' ? undefined : navigateTo('/no-access')
	}

	if (to.path === '/login' || to.path === '/no-access') {
		return navigateTo('/')
	}
})
