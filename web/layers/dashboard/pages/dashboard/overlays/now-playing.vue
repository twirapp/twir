<script setup lang="ts">
import { NowPlaying, Preset } from '@twir/frontend-now-playing'
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from 'reka-ui'
import { computed, ref, watch } from 'vue'
import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useNowPlayingOverlayApi } from '~~/layers/dashboard/api/overlays/now-playing'
import { useTheme } from '~~/layers/dashboard/composables/use-theme'
import NowPlayingForm from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/now-playing-form.vue'
import {
	defaultSettings,
	useNowPlayingForm,
} from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-form'
import { useNowPlayingPreviewTrack } from '~~/layers/dashboard/pages/dashboard/overlays/now-playing/use-now-playing-track'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { theme } = useTheme()
const { t } = useI18n()

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const nowPlayingOverlayManager = useNowPlayingOverlayApi()
const creator = nowPlayingOverlayManager.useNowPlayingCreate()

const { track: nowPlayingTrack, isSomeSongIntegrationEnabled } = useNowPlayingPreviewTrack()

const { data: settings, setData } = useNowPlayingForm()

const { data: entities } = nowPlayingOverlayManager.useNowPlayingQuery()

const openedTab = ref<string>()

async function handleAdd() {
	const input = { ...defaultSettings }
	// @ts-expect-error: Create input excludes the persisted channel ID.
	delete input.channelId
	// @ts-expect-error: Create input excludes the persisted overlay ID.
	delete input.id

	await creator.executeMutation({
		input,
	})
}
const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.nowPlayingOverlays.length ?? 0) < 5
})

watch(
	entities,
	(newValue, oldValue) => {
		if (newValue?.nowPlayingOverlays.length === oldValue?.nowPlayingOverlays.length) {
			return
		}

		if (!entities.value?.nowPlayingOverlays.at(0)) {
			openedTab.value = undefined
			return
		}

		openedTab.value = entities.value.nowPlayingOverlays.at(0)!.id
	},
	{ immediate: true }
)

watch(openedTab, async (v) => {
	const entity = entities.value?.nowPlayingOverlays.find((s) => s.id === v)
	if (!entity) return
	setData(entity)
})
</script>

<template>
	<div class="flex flex-col gap-3">
		<div>
			<NowPlaying
				:settings="settings ?? { preset: Preset.TRANSPARENT }"
				:track="nowPlayingTrack"
			/>
		</div>
		<Separator />
		<Alert
			v-if="!isSomeSongIntegrationEnabled"
			variant="destructive"
		>
			<Icon
				name="lucide:alert-triangle"
				class="h-4 w-4"
			/>
			<AlertTitle>No enabled song integrations!</AlertTitle>
			<AlertDescription>
				Connect Spotify, Last.fm or VK in
				<NuxtLink
					to="/dashboard/integrations"
					class="text-primary hover:underline"
				>
					{{ t('sidebar.integrations') }}
				</NuxtLink>
				to use this overlay
			</AlertDescription>
		</Alert>
		<TabsRoot
			v-else
			:default-value="0"
			orientation="vertical"
			class="min-h-[45dvh]"
			@update:model-value="
				(e) => {
					const overlay = entities?.nowPlayingOverlays[Number(e)]
					openedTab = overlay?.id
				}
			"
		>
			<TabsList
				aria-label="tabs example"
				class="-mb-px flex flex-wrap items-center overflow-x-auto"
			>
				<Button
					size="sm"
					variant="secondary"
					class="mr-1"
					:disabled="!addable"
					@click="handleAdd"
				>
					<Icon name="lucide:plus" />
				</Button>
				<TabsTrigger
					v-for="(overlay, index) of entities?.nowPlayingOverlays"
					:key="overlay.id"
					class="relative z-10 flex px-3 py-4 text-sm font-medium whitespace-nowrap transition-colors before:absolute before:top-2 before:left-0 before:-z-10 before:block before:h-9 before:w-full before:rounded-md before:transition-colors before:content-[''] hover:text-white hover:before:bg-zinc-800 data-disabled:cursor-not-allowed data-disabled:text-zinc-400 data-[state=active]:after:absolute data-[state=active]:after:right-2 data-[state=active]:after:bottom-0 data-[state=active]:after:left-2 data-[state=active]:after:block data-[state=active]:after:h-0 data-[state=active]:after:rounded-t-sm data-[state=active]:after:border-b-2 data-[state=active]:after:content-['']"
					:class="[
						theme === 'dark'
							? 'data-[state=active]:after:border-white'
							: 'data-[state=active]:after:border-zinc-800',
					]"
					:value="index"
				>
					#{{ index + 1 }} {{ overlay.preset }}
				</TabsTrigger>
			</TabsList>
			<Alert
				v-if="!entities?.nowPlayingOverlays.length"
				class="mt-2"
			>
				<AlertTitle>No overlays!</AlertTitle>
				<AlertDescription> Create new overlay for edit settings </AlertDescription>
			</Alert>
			<TabsContent
				v-for="(overlay, index) of entities?.nowPlayingOverlays"
				:key="overlay.id"
				class="mt-2"
				:value="index"
			>
				<NowPlayingForm />
			</TabsContent>
		</TabsRoot>
	</div>
</template>

<style scoped>
.iframe {
	border: 1px solid hsl(var(--border));
	border-radius: 8px;
	padding: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	background-position: center;
	background-repeat: no-repeat;
	background-size: cover;
}
</style>
