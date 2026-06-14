import { profileQuery, userAccessFlagChecker } from '../api/auth'

export default defineNuxtRouteMiddleware(async (to) => {
	if (to.path.startsWith('/dashboard/popup')) return

	const urqlClient = useUrqlClient()
	const { data } = await urqlClient.query(profileQuery, {}).toPromise()

	if (!data?.authenticatedUser) {
		return navigateTo('/', { replace: true })
	}

	if (to.meta.adminOnly && !data.authenticatedUser.isBotAdmin) {
		return navigateTo('/dashboard/forbidden', { replace: true })
	}

	if (to.meta.neededPermission) {
		const hasAccess = await userAccessFlagChecker(to.meta.neededPermission as any)
		if (!hasAccess) {
			return navigateTo('/dashboard/forbidden', { replace: true })
		}
	}
})
