import { useTTSOverlayApi } from '#layers/dashboard/api/overlays-tts'
import { computed } from 'vue'

interface VoiceInfo {
	country: string
	gender: string
	lang: string
	name: string
	no: number
}

const countriesMapping: Record<string, string> = {
	ru: '🇷🇺',
	mk: '🇲🇰',
	uk: '🇺🇦',
	ka: '🇬🇪',
	ky: '🇰🇬',
	en: '🇺🇸',
	pt: '🇵🇹',
	eo: '🇺🇳',
	sq: '🇦🇱',
	cs: '🇨🇿',
	pl: '🇵🇱',
	br: '🇧🇷',
}

export interface VoiceOption {
	value: string
	label: string
	lang: string
}

export function useTTSVoices() {
	const api = useTTSOverlayApi()
	const { data, fetching } = api.useQueryTTSGetInfo()

	const voicesInfo = computed<Record<string, VoiceInfo>>(() => {
		if (!data.value?.overlaysTTSGetInfo?.voicesInfo) return {}

		const result: Record<string, VoiceInfo> = {}
		for (const item of data.value.overlaysTTSGetInfo.voicesInfo) {
			result[item.key] = {
				country: item.info.country,
				gender: item.info.gender,
				lang: item.info.lang,
				name: item.info.name,
				no: item.info.no,
			}
		}
		return result
	})

	const voices = computed<VoiceOption[]>(() => {
		const result: VoiceOption[] = []

		for (const [voiceKey, voice] of Object.entries(voicesInfo.value)) {
			const flag = countriesMapping[voice.lang] || ''
			result.push({
				value: voiceKey,
				label: `${flag} ${voice.name} (${voice.gender})`,
				lang: voice.lang,
			})
		}

		return result.sort((a, b) => a.label.localeCompare(b.label))
	})

	const voicesByLanguage = computed(() => {
		const grouped: Record<string, VoiceOption[]> = {}

		for (const voice of voices.value) {
			if (!grouped[voice.lang]) {
				grouped[voice.lang] = []
			}
			grouped[voice.lang].push(voice)
		}

		return grouped
	})

	return {
		voices,
		voicesByLanguage,
		voicesInfo,
		isLoading: fetching,
	}
}
