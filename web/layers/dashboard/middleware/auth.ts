import { profileQuery, userAccessFlagChecker } from '../api/auth'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'
import { parseMcpConsentAttempt } from '~/utils/mcp-consent.js'

export function isChannelRolePermission(value: unknown): value is ChannelRolePermissionEnum {
	return typeof value === 'string' && Object.values(ChannelRolePermissionEnum).some((permission) => permission === value)
}

export default defineNuxtRouteMiddleware(async (to) => {
	if (to.path.startsWith('/dashboard/popup')) return

	const localePath = useLocalePath()
	const urqlClient = useUrqlClient()
	const { data } = await urqlClient.query(profileQuery, {}).toPromise()

	if (!data?.authenticatedUser) {
		const attempt = to.path === '/dashboard/mcp/authorize'
			? parseMcpConsentAttempt(to.query.attempt)
			: null
		if (attempt !== null) {
			return navigateTo({ path: '/', query: { mcp_attempt: attempt } }, { replace: true })
		}

		return navigateTo('/', { replace: true })
	}

	if (to.meta.adminOnly && !data.authenticatedUser.isBotAdmin) {
		return navigateTo(localePath('/dashboard/forbidden'), { replace: true })
	}

	if (to.meta.neededPermission !== undefined) {
		if (!isChannelRolePermission(to.meta.neededPermission)) {
			return navigateTo(localePath('/dashboard/forbidden'), { replace: true })
		}

		const hasAccess = await userAccessFlagChecker(to.meta.neededPermission)
		if (!hasAccess) {
			return navigateTo(localePath('/dashboard/forbidden'), { replace: true })
		}
	}
})
