import { createGlobalState } from '@vueuse/core'
import { ref, toRaw } from 'vue'

import type { NowPlayingOverlay } from '~/gql/graphql'

import { NowPlayingOverlayPreset } from '~/gql/graphql'

export const defaultSettings: NowPlayingOverlay = {
	id: '',
	preset: NowPlayingOverlayPreset.Transparent,
	backgroundColor: 'rgba(0, 0, 0, 0)',
	channelId: '',
	fontFamily: 'inter',
	fontWeight: 400,
	showImage: true,
	hideTimeout: 0,
}

export const useNowPlayingForm = createGlobalState(() => {
	const data = ref<NowPlayingOverlay>(structuredClone(defaultSettings))

	function setData(d: NowPlayingOverlay) {
		data.value = structuredClone(toRaw(d))
	}

	function getDefaultSettings() {
		return structuredClone(defaultSettings)
	}

	return {
		data,
		setData,
		getDefaultSettings,
	}
})
