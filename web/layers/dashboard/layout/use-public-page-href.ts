import { computed } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth'

export function usePublicPageHref() {
	const { data: profileData } = useProfile()
	const selectedDashboardTwitchUser = computed(() => {
		return profileData.value?.availableDashboards.find(
			(d) => d.id === profileData.value?.selectedDashboardId
		)?.twitchProfile
	})

	return computed(() => {
		const selectedDashboardLogin = selectedDashboardTwitchUser.value?.login
		if (!selectedDashboardLogin) {
			return null
		}

		const origin = typeof window !== 'undefined'
			? window.location.origin
			: ''
		return `${origin}/p/${selectedDashboardLogin}`
	})
}
