<script setup lang="ts">
import type { AcceptableValue } from 'reka-ui'
import { Platform } from '~/gql/graphql.js'
import { PLATFORM_META } from '@/utils/platforms'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { useChannelPlatformsApi } from '../../channel-platforms/api'

interface Props { zoom: number }
interface StreamBackgroundBinding { readonly platform: Platform; readonly platformLogin: string }
interface StreamBackgroundPreference { readonly platform: Platform | null; readonly enabled: boolean }
type StreamPlayerUrlFactory = (login: string) => string

const props = defineProps<Props>()
const { t } = useI18n()
const { data: channelPlatforms } = useChannelPlatformsApi().useQuery()
const streamBackgroundPreference = useLocalStorage<StreamBackgroundPreference>(
	'overlayBuilder.streamBackground',
	{ platform: null, enabled: false },
	{ initOnMounted: true },
)
const streamPlayerUrlFactories = {
	[Platform.Twitch]: (login: string) => `https://player.twitch.tv/?channel=${encodeURIComponent(login)}&parent=${encodeURIComponent(window.location.hostname)}&autoplay=true&muted=true`,
	[Platform.Kick]: (login: string) => `https://player.kick.com/${encodeURIComponent(login)}?autoplay=true&muted=true`,
	[Platform.VkVideoLive]: null,
	[Platform.Youtube]: null,
} satisfies Record<Platform, StreamPlayerUrlFactory | null>
const enabledBindings = computed<StreamBackgroundBinding[]>(() => {
	const bindings: StreamBackgroundBinding[] = []
	for (const binding of channelPlatforms.value?.channelPlatformBindings ?? []) {
		if (binding.enabled && binding.platformLogin) bindings.push({ platform: binding.platform, platformLogin: binding.platformLogin })
	}
	return bindings
})
const selectedBinding = computed(() => enabledBindings.value.find((binding) => binding.platform === streamBackgroundPreference.value.platform) ?? enabledBindings.value[0] ?? null)
const streamPreviewSrc = computed(() => {
	const binding = selectedBinding.value
	if (!import.meta.client || !streamBackgroundPreference.value.enabled || !binding) return null
	return streamPlayerUrlFactories[binding.platform]?.(binding.platformLogin) ?? null
})
const showUnsupportedPreview = computed(() => streamBackgroundPreference.value.enabled && selectedBinding.value !== null && streamPreviewSrc.value === null)
const controlStyle = computed(() => ({ transform: `scale(${1 / props.zoom})`, transformOrigin: 'top right', zIndex: 10001 }))
function updatePlatform(value: AcceptableValue) {
	if (typeof value !== 'string') return
	const binding = enabledBindings.value.find((item) => item.platform === value)
	if (!binding) return
	streamBackgroundPreference.value = { ...streamBackgroundPreference.value, platform: binding.platform }
}
function updateEnabled(enabled: boolean) {
	streamBackgroundPreference.value = {
		...streamBackgroundPreference.value,
		platform: selectedBinding.value?.platform ?? streamBackgroundPreference.value.platform,
		enabled,
	}
}
</script>

<template>
	<div v-if="streamPreviewSrc || showUnsupportedPreview" class="pointer-events-none absolute inset-0 overflow-hidden" :style="{ zIndex: 0 }">
		<iframe v-if="streamPreviewSrc" class="pointer-events-none size-full border-0" :src="streamPreviewSrc" :title="t('overlayBuilder.streamBackground.previewTitle')" allow="autoplay" />
		<div v-else class="flex size-full items-center justify-center bg-background/80 p-4 text-center text-xs text-muted-foreground">{{ t('overlayBuilder.streamBackground.previewNotSupported') }}</div>
	</div>
	<div v-if="enabledBindings.length" class="absolute right-4 top-4" :style="controlStyle">
		<div class="flex items-center gap-2 rounded-md border bg-background p-2 shadow-sm">
			<Select v-if="enabledBindings.length > 1" :model-value="selectedBinding?.platform" @update:model-value="updatePlatform">
				<SelectTrigger class="h-7 w-32 text-xs"><SelectValue /></SelectTrigger>
				<SelectContent><SelectItem v-for="binding in enabledBindings" :key="binding.platform" :value="binding.platform"><span class="flex items-center gap-2"><Icon :name="PLATFORM_META[binding.platform].icon" :class="PLATFORM_META[binding.platform].colorClass" />{{ PLATFORM_META[binding.platform].label }}</span></SelectItem></SelectContent>
			</Select>
			<Label for="stream-background" class="cursor-pointer whitespace-nowrap text-xs">{{ t('overlayBuilder.streamBackground.toggle') }}</Label>
			<Switch id="stream-background" :model-value="streamBackgroundPreference.enabled" @update:model-value="updateEnabled" />
		</div>
	</div>
</template>
