import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'

import type { KappagenSettings } from '@/types.js'
import type { KappagenConfig } from '@twirapp/kappagen/types'

export const useKappagenSettings = createGlobalState(() => {
	const kappagenSettings = ref<Required<KappagenConfig>>({
		max: 0,
		time: 5,
		queue: 0,
		cube: {
			speed: 6,
		},
		animation: {
			fade: {
				in: 8,
				out: 9,
			},
			zoom: {
				in: 17,
				out: 8,
			},
		},
		in: {
			fade: true,
			zoom: true,
		},
		out: {
			fade: true,
			zoom: true,
		},
		size: {
			min: 1,
			max: 256,
			ratio: {
				normal: 1 / 12,
				small: 1 / 24,
			},
		},
	})

	const overlaySettings = ref<KappagenSettings>()

	function updateOverlaySettings(settings: KappagenSettings): void {
		overlaySettings.value = settings

		if (settings.emotes) {
			kappagenSettings.value.max = settings.emotes.max
			kappagenSettings.value.time = settings.emotes.time
			kappagenSettings.value.queue = settings.emotes.queue
		}

		if (settings.size) {
			kappagenSettings.value.size = {
				min: settings.size.min,
				max: settings.size.max,
				ratio: {
					normal: settings.size.ratioNormal,
					small: settings.size.ratioSmall,
				},
			}
		}

		if (settings.cube) {
			kappagenSettings.value.cube = {
				speed: settings.cube.speed,
			}
		}

		if (settings.animation) {
			kappagenSettings.value.in = {
				fade: settings.animation.fadeIn,
				zoom: settings.animation.zoomIn,
			}

			kappagenSettings.value.out = {
				fade: settings.animation.fadeOut,
				zoom: settings.animation.zoomOut,
			}
		}
	}

	return {
		overlaySettings,
		kappagenSettings,
		updateSettings: updateOverlaySettings,
	}
})
