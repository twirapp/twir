import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api'

export function useAlertsInformation() {
	const { t } = useI18n()

	const { data: profile } = useProfile()
	const overlayLink = computed(() => {
		return `${window.location.origin}/overlays/${profile.value?.apiKey}/alerts`
	})

	const isShowOverlayLink = ref(false)
	function toggleShowOverlayLink() {
		isShowOverlayLink.value = !isShowOverlayLink.value
	}

	const copyButtonText = ref(t('alerts.copyOverlayLink'))
	async function copyOverlayLink() {
		try {
			await navigator.clipboard.writeText(overlayLink.value)
			copyButtonText.value = t('sharedTexts.copied')
			setTimeout(() => (copyButtonText.value = t('alerts.copyOverlayLink')), 1000)
		} catch (error) {
			console.error(error)
		}
	}

	return {
		overlayLink,
		isShowOverlayLink,
		toggleShowOverlayLink,

		copyButtonText,
		copyOverlayLink,
	}
}
