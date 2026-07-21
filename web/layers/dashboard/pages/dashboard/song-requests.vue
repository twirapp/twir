<script setup lang="ts">
import { toast } from 'vue-sonner'
import { useSongRequestsApi } from '~~/layers/dashboard/api/song-requests.js'
import OverlaySettingsModal from '~~/layers/dashboard/components/songRequests/overlay-settings.vue'
import Player from '~~/layers/dashboard/components/songRequests/player.vue'
import VideosQueue from '~~/layers/dashboard/components/songRequests/queue.vue'
import SettingsModal from '~~/layers/dashboard/components/songRequests/settings.vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()
const isSettingsModalOpened = ref(false)
const isOverlaySettingsModalOpened = ref(false)
const openSettingsModal = () => (isSettingsModalOpened.value = true)

const youtubeModuleManager = useSongRequestsApi()
const youtubeModuleData = youtubeModuleManager.useSongRequestQuery()

const channelApiKey = computed(() => {
	return youtubeModuleData.data.value?.songRequests?.channelApiKey ?? ''
})

const requestUrl = useRequestURL()

const widgetLink = computed(() => {
	if (!channelApiKey.value) return ''
	return `${requestUrl.origin}/w/${channelApiKey.value}/song-requests`
})

const overlayLink = computed(() => {
	if (!channelApiKey.value) return ''
	return `${requestUrl.origin}/o/${channelApiKey.value}/song-requests`
})

const showLinks = ref(false)

function copyLink(link: string, label: string) {
	if (!link) {
		toast.error(t('songRequests.links.copyError'), { duration: 2500 })
		return
	}

	navigator.clipboard
		.writeText(link)
		.then(() => {
			toast.success(t('songRequests.links.copied', { label }), { duration: 3000 })
		})
		.catch(() => {
			toast.error(t('songRequests.links.copyError'), { duration: 2500 })
		})
}
</script>

<template>
	<Card class="mb-4">
		<CardHeader>
			<CardTitle>{{ t('songRequests.links.title') }}</CardTitle>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<div
				v-if="!youtubeModuleData.fetching.value && !channelApiKey"
				class="text-muted-foreground text-sm"
			>
				{{ t('songRequests.links.notConfigured') }}
			</div>
			<template v-else>
				<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
					<span class="min-w-24 text-sm font-medium">{{ t('songRequests.links.widget') }}:</span>
					<div class="relative flex-1">
						<Input
							:type="showLinks ? 'text' : 'password'"
							:model-value="widgetLink"
							readonly
							class="pr-24 font-mono text-sm"
						/>
						<div class="absolute top-1/2 right-1 flex -translate-y-1/2 gap-1">
							<Button
								variant="ghost"
								size="sm"
								class="h-7 px-2 text-xs"
								@click="showLinks = !showLinks"
							>
								{{ showLinks ? t('sharedButtons.hide') : t('sharedButtons.show') }}
							</Button>
							<Button
								variant="outline"
								size="sm"
								class="h-7 px-2 text-xs"
								@click="copyLink(widgetLink, t('songRequests.links.widget'))"
							>
								{{ t('sharedButtons.copy') }}
							</Button>
						</div>
					</div>
				</div>
				<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-3">
					<span class="min-w-24 text-sm font-medium">{{ t('songRequests.links.overlay') }}:</span>
					<div class="relative flex-1">
						<Input
							:type="showLinks ? 'text' : 'password'"
							:model-value="overlayLink"
							readonly
							class="pr-36 font-mono text-sm"
						/>
						<div class="absolute top-1/2 right-1 flex -translate-y-1/2 gap-1">
							<Button
								variant="ghost"
								size="sm"
								class="h-7 px-2 text-xs"
								@click="showLinks = !showLinks"
							>
								{{ showLinks ? t('sharedButtons.hide') : t('sharedButtons.show') }}
							</Button>
							<Button
								variant="outline"
								size="sm"
								class="h-7 px-2 text-xs"
								@click="copyLink(overlayLink, t('songRequests.links.overlay'))"
							>
								{{ t('sharedButtons.copy') }}
							</Button>
							<Button
								type="button"
								variant="outline"
								size="icon"
								class="size-7"
								:aria-label="t('songRequests.overlaySettings.button')"
								:title="t('songRequests.overlaySettings.button')"
								@click="isOverlaySettingsModalOpened = true"
							>
								<Icon name="lucide:settings" />
							</Button>
						</div>
					</div>
				</div>
			</template>
		</CardContent>
	</Card>

	<div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
		<div class="lg:col-span-1">
			<Player
				v-if="!youtubeModuleData.fetching.value"
				:open-settings-modal="openSettingsModal"
			/>
		</div>

		<div class="lg:col-span-2">
			<VideosQueue />
		</div>
	</div>

	<SettingsModal v-model:open="isSettingsModalOpened" />
	<OverlaySettingsModal v-model="isOverlaySettingsModalOpened" />
</template>
