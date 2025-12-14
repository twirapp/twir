<script setup lang="ts">
import { NowPlaying, Preset } from '@twir/frontend-now-playing'
import { useSubscription } from '@urql/vue'
import { PlusIcon } from 'lucide-vue-next'
import { NA, NResult, useThemeVars } from 'naive-ui'
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from 'radix-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useNowPlayingOverlayApi, useProfile, useUserAccessFlagChecker } from '@/api'
import { useIntegrationsPageData } from '@/api/integrations/integrations-page.ts'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useTheme } from '@/composables/use-theme.ts'
import { graphql } from '@/gql'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import NowPlayingForm from '@/pages/overlays/now-playing/now-playing-form.vue'
import {
	defaultSettings,
	useNowPlayingForm,
} from '@/pages/overlays/now-playing/use-now-playing-form'

const { theme } = useTheme()
const themeVars = useThemeVars()
const { t } = useI18n()
const profile = useProfile()

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const nowPlayingOverlayManager = useNowPlayingOverlayApi()
const creator = nowPlayingOverlayManager.useNowPlayingCreate()

const integrationsPage = useIntegrationsPageData()

const isSomeSongIntegrationEnabled = computed(() => {
	return (
		integrationsPage.spotifyData.value?.userName ||
		integrationsPage.lastfmData.value?.userName ||
		integrationsPage.vkData.value?.userName
	)
})

const { data: settings, setData } = useNowPlayingForm()

const { data: entities } = nowPlayingOverlayManager.useNowPlayingQuery()

const openedTab = ref<string>()

async function handleAdd() {
	const input = { ...defaultSettings }
	// eslint-disable-next-line ts/ban-ts-comment
	// @ts-expect-error
	delete input.channelId
	// eslint-disable-next-line ts/ban-ts-comment
	// @ts-expect-error
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

const currentTrackPaused = computed(() => {
	return !isSomeSongIntegrationEnabled.value || !profile.data.value?.apiKey
})

const { data: currentTrackSub } = useSubscription({
	query: graphql(`
		subscription NowPlayingOverlayNowPlaying($apiKey: String!) {
			nowPlayingCurrentTrack(apiKey: $apiKey) {
				title
				artist
				imageUrl
			}
		}
	`),
	get variables() {
		return {
			apiKey: profile.data.value!.apiKey!,
		}
	},
	pause: currentTrackPaused,
})

const defaultNowPlayingTrack = {
	imageUrl: 'https://i.scdn.co/image/ab67616d0000b273e7fbc0883149094912559f2c',
	artist: 'Slipknot',
	title: 'Psychosocial',
}

const nowPlayingTrack = computed(() => {
	if (currentTrackSub.value?.nowPlayingCurrentTrack) {
		return {
			imageUrl: currentTrackSub.value.nowPlayingCurrentTrack.imageUrl,
			artist: currentTrackSub.value.nowPlayingCurrentTrack.artist,
			title: currentTrackSub.value.nowPlayingCurrentTrack.title,
		}
	}

	return defaultNowPlayingTrack
})
</script>

<template>
	<div class="flex flex-col gap-3">
		<div>
			<NowPlaying :settings="settings ?? { preset: Preset.TRANSPARENT }" :track="nowPlayingTrack" />
		</div>
		<Separator />
		<div>
			<NResult
				v-if="!isSomeSongIntegrationEnabled"
				status="warning"
				title="No enabled song integrations!"
			>
				<template #footer>
					Connect Spotify, Last.fm or VK in
					<router-link :to="{ name: 'Integrations' }" #="{ navigate, href }" custom>
						<NA :href="href" @click="navigate">
							{{ t('sidebar.integrations') }}
						</NA>
					</router-link>
					to use this overlay
				</template>
			</NResult>
			<TabsRoot
				v-else
				:default-value="0"
				orientation="vertical"
				class="min-h-[45dvh]"
				@update:model-value="(e) => (openedTab = entities?.nowPlayingOverlays[e].id)"
			>
				<TabsList
					aria-label="tabs example"
					class="flex flex-wrap items-center overflow-x-auto -mb-px"
				>
					<Button
						size="sm"
						variant="secondary"
						class="mr-1"
						:disabled="!addable"
						@click="handleAdd"
					>
						<PlusIcon />
					</Button>
					<TabsTrigger
						v-for="(overlay, index) of entities?.nowPlayingOverlays"
						:key="overlay.id"
						class="tabs-trigger data-disabled:cursor-not-allowed data-disabled:text-zinc-400"
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
				<Alert v-if="!entities?.nowPlayingOverlays.length" class="mt-2">
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
	</div>
</template>

<style scoped>
@reference '@/assets/index.css';
@import '../styles.css';

.iframe {
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
	padding: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
	background-position: center;
	background-repeat: no-repeat;
	background-size: cover;
}

.tabs-trigger {
	@apply relative z-1 flex whitespace-nowrap px-3 py-4 text-sm  transition-colors before:absolute before:left-0 before:top-2 before:-z-1 before:block before:h-9 before:w-full before:rounded-md before:transition-colors before:content-[''] hover:text-white hover:before:bg-zinc-800 data-[state=active]:after:absolute data-[state=active]:after:bottom-0 data-[state=active]:after:left-2 data-[state=active]:after:right-2 data-[state=active]:after:block data-[state=active]:after:h-0 data-[state=active]:after:border-b-2 data-[state=active]:after:content-[''] data-[state=active]:after:rounded-t-sm font-medium;
}
</style>
