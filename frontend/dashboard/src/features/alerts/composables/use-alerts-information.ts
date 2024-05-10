import { computed, ref } from 'vue'

import { useProfile } from '@/api'

export function useAlertsInformation() {
	const { data: profile } = useProfile()
	const overlayLink = computed(() => {
		return `${window.location.origin}/overlays/${profile.value?.apiKey}/alerts`
	})

	const isShowOverlayLink = ref(false)
	function toggleShowOverlayLink() {
		isShowOverlayLink.value = !isShowOverlayLink.value
	}

	async function copyOverlayLink() {
		try {
			await navigator.clipboard.writeText(overlayLink.value)
		} catch (error) {
			console.error(error)
		}
	}

	return {
		overlayLink,
		isShowOverlayLink,
		toggleShowOverlayLink,
		copyOverlayLink,
	}
}
