import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/auth.js'
import { useToast } from '@/components/ui/toast'

export function useCopyOverlayLink(overlayPath: string) {
	const { data: profile } = useProfile()
	const { t } = useI18n()
	const { toast } = useToast()

	const overlayLink = computed(() => {
		return `${window.location.origin}/overlays/${profile.value?.apiKey}/${overlayPath}`
	})

	const copyOverlayLink = (query?: Record<string, string>) => {
		const url = new URL(overlayLink.value)
		if (query) {
			for (const [key, value] of Object.entries(query)) {
				url.searchParams.set(key, value)
			}
		}

		navigator.clipboard.writeText(url.toString())

		toast({
			title: t('overlays.copied'),
			duration: 5000,
		})
	}

	return {
		overlayLink,
		copyOverlayLink,
	}
}
