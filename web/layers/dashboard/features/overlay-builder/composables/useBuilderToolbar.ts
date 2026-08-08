import { type MaybeRefOrGetter, computed, toValue } from 'vue'
import { toast } from 'vue-sonner'

import { useProfile } from '~~/layers/dashboard/api/auth.js'

export function useBuilderToolbar(overlayId: MaybeRefOrGetter<string | undefined>) {
	const { t } = useI18n()
	const router = useRouter()
	const localePath = useLocalePath()
	const { data: profile } = useProfile()
	const requestUrl = useRequestURL()

	const selectedDashboardUser = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const overlayApiKey = computed(() => {
		return selectedDashboardUser.value?.channelApiKey || profile.value?.channelApiKey || ''
	})

	const formatZoom = (zoom: number) => `${Math.round(zoom * 100)}%`

	function goBack() {
		router.push(localePath('/dashboard/overlays'))
	}

	function copyOverlayLink() {
		const id = toValue(overlayId)
		if (!id) return

		if (!overlayApiKey.value) {
			toast.error('No API key found')
			return
		}

		const overlayUrl = `${requestUrl.origin}/overlays/${overlayApiKey.value}/registry/overlays/${id}`
		navigator.clipboard
			.writeText(overlayUrl)
			.then(() => toast.success(t('sharedTexts.copied')))
			.catch(() => toast.error('Failed to copy link'))
	}

	return { formatZoom, goBack, copyOverlayLink }
}
