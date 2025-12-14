import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { Settings } from '@twir/frontend-valorant-stats'

import { useProfile } from '@/api/auth.ts'
import { toast } from 'vue-sonner'

export const useValorantStats = createGlobalState(() => {
	const { data: profile } = useProfile()

	const settings = ref<Required<Settings>>({
		backgroundColor: '#07090e',
		textColor: '#f2f2f2',
		primaryTextColor: '#B9B4B4',
		winColor: '#00FFE3',
		loseColor: '#FF7986',

		disabledPeakRR: false,
		disabledLeaderboardPlace: false,
		disabledPeakRankIcon: false,
		disabledBorder: false,
		disabledWinLose: false,
		disabledProgress: false,
		disabledGlowEffect: false,
		disabledTwentyLastMatches: false,
	})

	function setSettings(data: Required<Settings>) {
		settings.value = data
	}

	function copyOverlayLink() {
		const dashboard = profile.value?.availableDashboards.find(
			(d) => d.id === profile.value?.selectedDashboardId
		)

		const url = new URL(
			`${window.location.origin}/o/${dashboard?.apiKey ?? profile.value?.apiKey}/valorant-stats`
		)

		for (const [key, value] of Object.entries(settings.value)) {
			url.searchParams.set(key, value as string)
		}

		navigator.clipboard.writeText(url.toString())
		toast.success('Overlay link copied to clipboard', {
			duration: 2500,
		})
	}

	return {
		settings,
		setSettings,
		copyOverlayLink,
	}
})
