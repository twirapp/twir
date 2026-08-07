import { createGlobalState } from '@vueuse/core'
import type { AcceptableValue } from 'reka-ui'

import { Platform } from '~/gql/graphql.js'

import { useChannelPlatformsApi } from '../../channel-platforms/api'

export interface StreamBackgroundBinding {
	readonly platform: Platform
	readonly platformLogin: string
}

interface StreamBackgroundPreference {
	readonly platform: Platform | null
	readonly enabled: boolean
}

type StreamPlayerUrlFactory = (login: string) => string

const streamPlayerUrlFactories = {
	[Platform.Twitch]: null,
	[Platform.Kick]: (login: string) =>
		`https://player.kick.com/${encodeURIComponent(login)}?autoplay=true&muted=true&controls=false`,
	[Platform.VkVideoLive]: null,
	[Platform.Youtube]: null,
} satisfies Record<Platform, StreamPlayerUrlFactory | null>

export const useStreamBackground = createGlobalState(() => {
	const { data: channelPlatforms } = useChannelPlatformsApi().useQuery()

	const preference = useLocalStorage<StreamBackgroundPreference>(
		'overlayBuilder.streamBackground',
		{ platform: null, enabled: false },
		{ initOnMounted: true }
	)

	const enabledBindings = computed<StreamBackgroundBinding[]>(() => {
		const bindings: StreamBackgroundBinding[] = []
		for (const binding of channelPlatforms.value?.channelPlatformBindings ?? []) {
			if (binding.enabled && binding.platformLogin) {
				bindings.push({ platform: binding.platform, platformLogin: binding.platformLogin })
			}
		}
		return bindings
	})

	const selectedBinding = computed(() => {
		return (
			enabledBindings.value.find((binding) => binding.platform === preference.value.platform) ??
			enabledBindings.value[0] ??
			null
		)
	})

	const streamPreviewSrc = computed(() => {
		const binding = selectedBinding.value
		if (!import.meta.client || !preference.value.enabled || !binding) return null
		return streamPlayerUrlFactories[binding.platform]?.(binding.platformLogin) ?? null
	})

	// Twitch preview is driven through the Twitch JS embed API instead of an iframe
	// (see useTwitchEmbedPlayer), so it is excluded from the unsupported state.
	const showUnsupportedPreview = computed(() => {
		return (
			preference.value.enabled &&
			selectedBinding.value !== null &&
			selectedBinding.value.platform !== Platform.Twitch &&
			streamPreviewSrc.value === null
		)
	})

	const isPreviewActive = computed(() => {
		if (!preference.value.enabled || selectedBinding.value === null) return false
		if (selectedBinding.value.platform === Platform.Twitch) return true
		return streamPreviewSrc.value !== null || showUnsupportedPreview.value
	})

	function updatePlatform(value: AcceptableValue) {
		if (typeof value !== 'string') return
		const binding = enabledBindings.value.find((item) => item.platform === value)
		if (!binding) return
		preference.value = { ...preference.value, platform: binding.platform }
	}

	function updateEnabled(enabled: boolean) {
		preference.value = {
			...preference.value,
			platform: selectedBinding.value?.platform ?? preference.value.platform,
			enabled,
		}
	}

	return {
		preference,
		enabledBindings,
		selectedBinding,
		streamPreviewSrc,
		showUnsupportedPreview,
		isPreviewActive,
		updatePlatform,
		updateEnabled,
	}
})
