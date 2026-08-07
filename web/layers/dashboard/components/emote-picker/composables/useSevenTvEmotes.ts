import * as z from 'zod'
import { type Ref, computed, ref, watch } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth.js'

const GLOBAL_EMOTE_SET_URL = 'https://7tv.io/v3/emote-sets/global'
const CHANNEL_EMOTE_SET_URL = 'https://7tv.io/v3/users/twitch'
const FALLBACK_EMOTE_URL = 'https://cdn.7tv.app/emote'

const emoteFileSchema = z.object({
	name: z.string(),
	format: z.string().optional(),
})

const emoteHostSchema = z.object({
	url: z.string().optional(),
	files: z.array(emoteFileSchema).optional(),
})

const emoteSchema = z.object({
	id: z.string(),
	name: z.string(),
	data: z.object({ host: emoteHostSchema.optional() }).optional(),
	host: emoteHostSchema.optional(),
})

const emoteSetSchema = z.object({ emotes: z.array(emoteSchema) })
const globalResponseSchema = emoteSetSchema
const channelResponseSchema = z.object({ emote_set: emoteSetSchema.nullable().optional() })

export type SevenTvEmote = z.infer<typeof emoteSchema>
export type SevenTvTab = 'global' | 'channel'
export type SevenTvLoadStatus = 'idle' | 'loading' | 'success' | 'error'

export interface SelectedSevenTvEmote {
	readonly url: string
	readonly name: string
	readonly provider: '7TV'
}

const tabValues: readonly SevenTvTab[] = ['global', 'channel']

type SelectEmote = (emote: SelectedSevenTvEmote) => void

export function useSevenTvEmotes(open: Ref<boolean>, selectEmote: SelectEmote) {
	const { t } = useI18n()
	const { data: profile } = useProfile()
	const activeTab = ref<SevenTvTab>('global')
	const searchQuery = ref('')

	const globalEmotes = ref<SevenTvEmote[]>([])
	const globalStatus = ref<SevenTvLoadStatus>('idle')
	const globalError = ref('')

	const channelEmotes = ref<SevenTvEmote[]>([])
	const channelStatus = ref<SevenTvLoadStatus>('idle')
	const channelError = ref('')
	const loadedChannelId = ref<string | null>(null)

	const selectedTwitchUserId = computed(() => {
		const currentProfile = profile.value
		const dashboard = currentProfile?.availableDashboards.find(
			(item) => item.id === currentProfile.selectedDashboardId
		)

		if (dashboard?.platform.toLowerCase() !== 'twitch') return null
		return dashboard.profile?.platformUserId ?? null
	})

	const currentEmotes = computed(() => (activeTab.value === 'global' ? globalEmotes.value : channelEmotes.value))
	const currentStatus = computed(() => (activeTab.value === 'global' ? globalStatus.value : channelStatus.value))
	const currentError = computed(() => (activeTab.value === 'global' ? globalError.value : channelError.value))
	const filteredEmotes = computed(() => {
		const query = searchQuery.value.trim().toLowerCase()
		if (!query) return currentEmotes.value
		return currentEmotes.value.filter((emote) => emote.name.toLowerCase().includes(query))
	})
	const emptyMessage = computed(() => {
		if (searchQuery.value.trim()) return t('emotePicker.empty.search')
		if (activeTab.value === 'channel' && !selectedTwitchUserId.value) {
			return t('emotePicker.empty.noChannel')
		}
		if (activeTab.value === 'channel') return t('emotePicker.empty.channelNoSet')
		return t('emotePicker.empty.noEmotes')
	})

	function getEmoteUrl(emote: SevenTvEmote): string {
		const host = emote.data?.host ?? emote.host
		const webpFile = host?.files?.find(
			(file) => file.name.toLowerCase() === '2x.webp' || (file.format?.toUpperCase() === 'WEBP' && file.name.startsWith('2x.'))
		)
		const baseUrl = host?.url
			? host.url.startsWith('//')
				? `https:${host.url}`
				: host.url
			: `${FALLBACK_EMOTE_URL}/${emote.id}`

		return `${baseUrl.replace(/\/$/, '')}/${webpFile?.name ?? '2x.webp'}`
	}

	async function loadGlobal() {
		if (globalStatus.value === 'loading' || globalStatus.value === 'success') return

		globalStatus.value = 'loading'
		globalError.value = ''
		try {
			const response = await fetch(GLOBAL_EMOTE_SET_URL)
			if (!response.ok) {
				globalError.value = t('emotePicker.errors.apiStatus', { status: response.status })
				globalStatus.value = 'error'
				return
			}
			const payload = globalResponseSchema.parse(await response.json())
			globalEmotes.value = payload.emotes
			globalStatus.value = 'success'
		} catch {
			globalError.value = t('emotePicker.errors.loadGlobal')
			globalStatus.value = 'error'
		}
	}

	async function loadChannel() {
		const channelId = selectedTwitchUserId.value
		if (!channelId) {
			channelEmotes.value = []
			channelStatus.value = 'success'
			loadedChannelId.value = null
			return
		}
		if (channelStatus.value === 'loading' || loadedChannelId.value === channelId) return

		channelStatus.value = 'loading'
		channelError.value = ''
		try {
			const response = await fetch(`${CHANNEL_EMOTE_SET_URL}/${channelId}`)
			if (response.status === 404) {
				channelEmotes.value = []
				loadedChannelId.value = channelId
				channelStatus.value = 'success'
				return
			}
			if (!response.ok) {
				channelError.value = t('emotePicker.errors.apiStatus', { status: response.status })
				channelStatus.value = 'error'
				return
			}
			const payload = channelResponseSchema.parse(await response.json())
			channelEmotes.value = payload.emote_set?.emotes ?? []
			loadedChannelId.value = channelId
			channelStatus.value = 'success'
		} catch {
			channelError.value = t('emotePicker.errors.loadChannel')
			channelStatus.value = 'error'
		}
	}

	function loadActiveTab() {
		if (activeTab.value === 'global') {
			void loadGlobal()
			return
		}
		void loadChannel()
	}

	function retryActiveTab() {
		if (activeTab.value === 'global') globalStatus.value = 'idle'
		else channelStatus.value = 'idle'
		loadActiveTab()
	}

	function handleSelectEmote(emote: SevenTvEmote) {
		selectEmote({ url: getEmoteUrl(emote), name: emote.name, provider: '7TV' })
		open.value = false
	}

	watch([open, activeTab], ([isOpen]) => {
		if (isOpen) loadActiveTab()
	})

	watch(selectedTwitchUserId, () => {
		if (loadedChannelId.value === selectedTwitchUserId.value) return
		channelEmotes.value = []
		channelStatus.value = 'idle'
		loadedChannelId.value = null
		if (open.value && activeTab.value === 'channel') loadChannel()
	})

	return {
		activeTab,
		searchQuery,
		currentEmotes,
		currentStatus,
		currentError,
		filteredEmotes,
		emptyMessage,
		tabValues,
		retryActiveTab,
		handleSelectEmote,
		getEmoteUrl,
	}
}
