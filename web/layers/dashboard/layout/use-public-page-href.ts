import { useProfile } from '~~/layers/dashboard/api/auth'

export function usePublicPageHref() {
	const { data: profileData } = useProfile()
	const requestUrl = useRequestURL()
	const selectedDashboard = computed(() => {
		return profileData.value?.availableDashboards.find(
			(d) => d.id === profileData.value?.selectedDashboardId
		)
	})

	return computed(() => {
		const dashboard = selectedDashboard.value
		const selectedDashboardLogin = dashboard?.profile?.platformLogin
		if (!selectedDashboardLogin) {
			return null
		}

		return `${requestUrl.origin}/p/${dashboard?.platform}/${selectedDashboardLogin}`
	})
}
