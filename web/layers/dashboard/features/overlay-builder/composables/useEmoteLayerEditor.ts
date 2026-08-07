import { type Ref, computed } from 'vue'

import type { Layer, LayerSettings } from '../types'

export interface SelectedEmote {
	readonly url: string
	readonly name: string
	readonly provider: '7TV'
}

type UpdateSettings = (updates: Partial<LayerSettings>) => void

export function useEmoteLayerEditor(layer: Ref<Layer>, updateSettings: UpdateSettings) {
	const emoteUrl = computed({
		get: () => layer.value.settings.emoteUrl,
		set: (value: string) => updateSettings({ emoteUrl: value, emoteName: '', emoteProvider: '' }),
	})

	function selectEmote(emote: SelectedEmote) {
		updateSettings({
			emoteUrl: emote.url,
			emoteName: emote.name,
			emoteProvider: emote.provider,
		})
	}

	return { emoteUrl, selectEmote }
}
