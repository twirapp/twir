import { computed } from 'vue'

import { useProfile } from '@/api'

export function usePublicPageHref() {
	const { data: profileData } = useProfile()

	return computed(() => {
		const selectedDashboardLogin = profileData.value?.selectedDashboardTwitchUser?.login
		if (!selectedDashboardLogin) {
			return null
		}

		return `${window.location.origin}/p/${selectedDashboardLogin}`
	})
}
