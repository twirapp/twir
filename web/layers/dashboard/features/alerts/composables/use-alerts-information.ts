import { computed, ref } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth'

export function useAlertsInformation() {
	const { data: profile } = useProfile()
	const requestUrl = useRequestURL()
	const selectedDashboardTwitchUser = computed(() => {
		return profile.value?.availableDashboards.find(
			(d) => d.id === profile.value?.selectedDashboardId
		)
	})

	const overlayApiKey = computed(() => {
		return selectedDashboardTwitchUser.value?.apiKey || profile.value?.apiKey || ''
	})

	const overlayLink = computed(() => {
		if (!overlayApiKey.value) return ''
		return `${requestUrl.origin}/overlays/${overlayApiKey.value}/alerts`
	})

	const isShowOverlayLink = ref(false)
	function toggleShowOverlayLink() {
		isShowOverlayLink.value = !isShowOverlayLink.value
	}

	const isCopied = ref(false)
	async function copyOverlayLink() {
		try {
			await navigator.clipboard.writeText(overlayLink.value)
			isCopied.value = true
			setTimeout(() => (isCopied.value = false), 1000)
		} catch (error) {
			console.error(error)
		}
	}

	return {
		overlayLink,
		isShowOverlayLink,
		toggleShowOverlayLink,
		isCopied,
		copyOverlayLink,
	}
}
