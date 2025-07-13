import { computed } from 'vue'

import { useProfile } from '@/api'

export function usePublicPageHref() {
	const { data: profileData } = useProfile()
	const selectedDashboardTwitchUser = computed(() => {
		return profileData.value?.availableDashboards.find((d) => d.id === profileData.value?.selectedDashboardId)
			?.twitchProfile
	})

	return computed(() => {
		const selectedDashboardLogin = selectedDashboardTwitchUser.value?.login
		if (!selectedDashboardLogin) {
			return null
		}

		return `${window.location.origin}/p/${selectedDashboardLogin}`
	})
}
