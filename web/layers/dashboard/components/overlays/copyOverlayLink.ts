import { computed } from 'vue'


import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { toast } from 'vue-sonner'

	export function useCopyOverlayLink(overlayPath: string, basePath: '/overlays' | '/o' = '/overlays') {
		const { data: profile } = useProfile()
		const { t } = useI18n()
		const requestUrl = useRequestURL()

	const selectedDashboardUser = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const overlayApiKey = computed(() => {
		return selectedDashboardUser.value?.channelApiKey || profile.value?.channelApiKey || ''
	})

	const overlayLink = computed(() => {
		if (!overlayApiKey.value) {
			return ''
		}

		return `${requestUrl.origin}${basePath}/${overlayApiKey.value}/${overlayPath}`
	})

		const canCopyOverlayLink = computed(() => Boolean(overlayLink.value))

	const copyOverlayLink = (query?: Record<string, string>) => {
		if (!overlayLink.value) {
			toast.error('Something went wrong at copying the overlay link', {
				duration: 2500,
			})
			return
		}

		const url = new URL(overlayLink.value)
		if (query) {
			for (const [key, value] of Object.entries(query)) {
				url.searchParams.set(key, value)
			}
		}

		navigator.clipboard.writeText(url.toString()).then(() => {
			toast.success(t('overlays.copied'), {
				duration: 5000,
			})
		}).catch(() => {
			toast.error('Failed to copy link to clipboard', {
				duration: 2500,
			})
		})
	}

		return {
			canCopyOverlayLink,
			overlayLink,
			copyOverlayLink,
		}
	}
