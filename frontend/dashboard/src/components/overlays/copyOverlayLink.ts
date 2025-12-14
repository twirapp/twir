import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/auth.js'
import { toast } from 'vue-sonner'

export function useCopyOverlayLink(overlayPath: string) {
	const { data: profile } = useProfile()
	const { t } = useI18n()

	const selectedDashboardUser = computed(() => {
		return profile.value?.availableDashboards.find(
			(dashboard) => dashboard.id === profile.value?.selectedDashboardId
		)
	})

	const overlayLink = computed(() => {
		if (!selectedDashboardUser.value?.apiKey) {
			return ''
		}

		return `${window.location.origin}/overlays/${selectedDashboardUser.value?.apiKey}/${overlayPath}`
	})

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

		navigator.clipboard.writeText(url.toString())

		toast.success(t('overlays.copied'), {
			duration: 5000,
		})
	}

	return {
		overlayLink,
		copyOverlayLink,
	}
}
