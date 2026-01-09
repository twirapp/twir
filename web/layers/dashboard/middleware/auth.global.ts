export default defineNuxtRouteMiddleware(async (to) => {
	// Only run on dashboard routes
	if (!to.path.startsWith('/dashboard')) return

	// Skip for popup routes (they may have different auth requirements)
	if (to.path.startsWith('/dashboard/popup')) return

	const authStore = useDashboardAuth()

	// Fetch user if not loaded - this will use SSR on server, client on browser
	if (!authStore.user && !authStore.isLoading) {
		const result = await authStore.fetchUser()

		// Wait a tick for the reactive data to update
		await nextTick()

		// Handle the case where the query returned but user is still null
		if (!authStore.user && !result.data?.value?.authenticatedUser) {
			return navigateTo('/', { external: true })
		}
	}

	// If still loading, wait for it to complete
	while (authStore.isLoading) {
		await new Promise((resolve) => setTimeout(resolve, 50))
	}

	// Redirect if no user or user is banned
	if (!authStore.user || authStore.user.isBanned) {
		return navigateTo('/', { external: true })
	}

	// Check admin-only routes
	if (to.meta.adminOnly && !authStore.user.isBotAdmin) {
		return navigateTo('/dashboard/forbidden')
	}

	// Check permissions
	if (to.meta.neededPermission) {
		const hasAccess = authStore.checkPermission(to.meta.neededPermission as any)
		if (!hasAccess) {
			return navigateTo('/dashboard/forbidden')
		}
	}
})
