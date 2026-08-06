<script setup lang="ts">
import * as z from 'zod'

import { useProfile } from '~~/layers/dashboard/api/auth.js'

import { Button } from '@/components/ui/button'
import { Dialog, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

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

type Emote = z.infer<typeof emoteSchema>
type Tab = 'global' | 'channel'
type LoadStatus = 'idle' | 'loading' | 'success' | 'error'

const tabValues: readonly Tab[] = ['global', 'channel']

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
	select: [emote: { url: string; name: string; provider: '7TV' }]
}>()

const { data: profile } = useProfile()
const activeTab = ref<Tab>('global')
const searchQuery = ref('')

const globalEmotes = ref<Emote[]>([])
const globalStatus = ref<LoadStatus>('idle')
const globalError = ref('')

const channelEmotes = ref<Emote[]>([])
const channelStatus = ref<LoadStatus>('idle')
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
	if (searchQuery.value.trim()) return 'По вашему запросу эмоции не найдены.'
	if (activeTab.value === 'channel' && !selectedTwitchUserId.value) {
		return 'Выберите Twitch-канал, чтобы увидеть его набор 7TV.'
	}
	if (activeTab.value === 'channel') return 'У текущего канала нет подключенного набора 7TV.'
	return 'В этом наборе пока нет эмоций.'
})

function getEmoteUrl(emote: Emote): string {
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
		if (!response.ok) throw new Error(`7TV returned ${response.status}`)
		const payload = globalResponseSchema.parse(await response.json())
		globalEmotes.value = payload.emotes
		globalStatus.value = 'success'
	} catch (error) {
		globalError.value = error instanceof Error ? error.message : 'Не удалось загрузить эмоции.'
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
		if (!response.ok) throw new Error(`7TV returned ${response.status}`)
		const payload = channelResponseSchema.parse(await response.json())
		channelEmotes.value = payload.emote_set?.emotes ?? []
		loadedChannelId.value = channelId
		channelStatus.value = 'success'
	} catch (error) {
		channelError.value = error instanceof Error ? error.message : 'Не удалось загрузить эмоции канала.'
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

function selectEmote(emote: Emote) {
	emit('select', { url: getEmoteUrl(emote), name: emote.name, provider: '7TV' })
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
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger v-if="$slots.trigger" as-child>
			<slot name="trigger" />
		</DialogTrigger>
		<DialogOrSheet class="max-h-[80dvh] max-w-3xl gap-0 overflow-hidden rounded-t-2xl p-0 sm:rounded-2xl">
			<div class="flex min-h-0 flex-col">
				<DialogHeader class="border-b px-6 py-4">
					<DialogTitle>Эмоции</DialogTitle>
					<DialogDescription class="sr-only">Выберите эмоцию из наборов 7TV.</DialogDescription>
				</DialogHeader>

				<Tabs v-model="activeTab" class="min-h-0 gap-0">
					<div class="sticky top-0 z-10 border-b bg-background/95 px-4 py-3 backdrop-blur">
						<TabsList class="grid h-9 w-full grid-cols-2 bg-muted/60">
							<TabsTrigger value="global">7TV</TabsTrigger>
							<TabsTrigger value="channel">7TV канал</TabsTrigger>
						</TabsList>
						<div class="relative mt-3">
							<Icon name="lucide:search" class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							<Input v-model="searchQuery" class="pl-9" placeholder="Поиск по имени" aria-label="Поиск эмоций" @keydown.stop />
						</div>
					</div>

					<TabsContent
						v-for="tab in tabValues"
						:key="tab"
						:value="tab"
						class="max-h-[60dvh] overflow-y-auto px-4 pb-4 pt-3"
					>
						<div v-if="currentStatus === 'loading'" class="grid grid-cols-4 gap-2 sm:grid-cols-6 md:grid-cols-8">
							<Skeleton v-for="index in 24" :key="index" class="aspect-square rounded-lg" />
						</div>
						<div v-else-if="currentStatus === 'error'" class="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
							<Icon name="lucide:triangle-alert" class="size-8 text-destructive" />
							<p class="max-w-sm text-sm text-muted-foreground">{{ currentError }}</p>
							<Button variant="outline" size="sm" @click="retryActiveTab">Повторить</Button>
						</div>
						<div v-else-if="filteredEmotes.length === 0" class="flex min-h-48 flex-col items-center justify-center gap-2 text-center">
							<Icon name="lucide:smile-plus" class="size-8 text-muted-foreground" />
							<p class="max-w-sm text-sm text-muted-foreground">{{ emptyMessage }}</p>
						</div>
						<div v-else class="grid grid-cols-4 gap-2 sm:grid-cols-6 md:grid-cols-8">
							<button
								v-for="emote in filteredEmotes"
								:key="emote.id"
								type="button"
								class="group flex aspect-square min-w-0 items-center justify-center rounded-lg border border-transparent bg-muted/30 p-2 transition-colors hover:border-border hover:bg-accent focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:scale-[0.98]"
								:title="emote.name"
								:aria-label="emote.name"
								@click="selectEmote(emote)"
							>
								<img
									:src="getEmoteUrl(emote)"
									:alt="emote.name"
									width="64"
									height="64"
									loading="lazy"
									class="size-full max-h-16 max-w-16 object-contain transition-transform group-hover:scale-110"
								/>
							</button>
						</div>
					</TabsContent>
				</Tabs>
			</div>
		</DialogOrSheet>
	</Dialog>
</template>
