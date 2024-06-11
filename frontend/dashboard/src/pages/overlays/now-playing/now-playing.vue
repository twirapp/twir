<script setup lang="ts">
import { NowPlaying, Preset } from '@twir/frontend-now-playing'
import { NA, NAlert, NResult, NTabPane, NTabs, useThemeVars } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import {
	useLastfmIntegration,
	useNowPlayingOverlayApi,
	useSpotifyIntegration,
	useUserAccessFlagChecker,
	useVKIntegration,
} from '@/api'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import NowPlayingForm from '@/pages/overlays/now-playing/now-playing-form.vue'
import {
	defaultSettings,
	useNowPlayingForm,
} from '@/pages/overlays/now-playing/use-now-playing-form'

const themeVars = useThemeVars()
const { t } = useI18n()
const { dialog } = useNaiveDiscrete()

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const nowPlayingOverlayManager = useNowPlayingOverlayApi()
const creator = nowPlayingOverlayManager.useNowPlayingCreate()
const deleter = nowPlayingOverlayManager.useNowPlayingDelete()

const { data: spotifyData } = useSpotifyIntegration().useData()
const { data: lastFmData } = useLastfmIntegration().useData()
const { data: vkData } = useVKIntegration().useData()

const isSomeSongIntegrationEnabled = computed(() => {
	return spotifyData.value?.userName || lastFmData.value?.userName || vkData.value?.userName
})

const {
	data: settings,
	setData,
} = useNowPlayingForm()

const {
	data: entities,
} = nowPlayingOverlayManager.useNowPlayingQuery()

const openedTab = ref<string>()

function resetTab() {
	if (!entities.value?.nowPlayingOverlays.at(0)) {
		openedTab.value = undefined
		return
	}

	openedTab.value = entities.value.nowPlayingOverlays.at(0)?.id
}

async function handleAdd() {
	await creator.executeMutation({
		input: defaultSettings,
	})
}

async function handleClose(id: string) {
	dialog.create({
		title: 'Delete preset',
		content: 'Are you sure you want to delete this preset?',
		positiveText: 'Delete',
		negativeText: 'Cancel',
		showIcon: false,
		onPositiveClick: async () => {
			const entity = entities.value?.nowPlayingOverlays.find(s => s.id === id)
			if (!entity?.id) return

			await deleter.executeMutation({
				id: entity.id,
			})
			resetTab()
		},
	})
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.nowPlayingOverlays.length ?? 0) < 5
})

watch(openedTab, async (v) => {
	const entity = entities.value?.nowPlayingOverlays.find(s => s.id === v)
	if (!entity) return
	setData(entity)
})

watch(entities, () => {
	resetTab()
}, { immediate: true })
</script>

<template>
	<div class="flex flex-col gap-3">
		<div>
			<NowPlaying
				:settings="settings ?? { preset: Preset.TRANSPARENT }"
				:track="{
					imageUrl: 'https://i.scdn.co/image/ab67616d0000b273e7fbc0883149094912559f2c',
					artist: 'Slipknot',
					title: 'Psychosocial',
				}"
			/>
		</div>
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
			<template v-if="isSomeSongIntegrationEnabled">
				<NTabs
					v-model:value="openedTab"
					type="card"
					:closable="userCanEditOverlays"
					:addable="addable"
					style="margin-top: 1rem;"
					tab-style="min-width: 80px;"
					@close="handleClose"
					@add="handleAdd"
				>
					<template #prefix>
						{{ t('overlays.chat.presets') }}
					</template>
					<template v-if="entities?.nowPlayingOverlays.length">
						<NTabPane
							v-for="(entity, entityIndex) in entities?.nowPlayingOverlays"
							:key="entity.id"
							:tab="`#${entityIndex + 1}`"
							:name="entity.id!"
						>
							<NowPlayingForm />
						</NTabPane>
					</template>
				</NTabs>
				<NAlert v-if="!entities?.nowPlayingOverlays.length" type="info" class="mt-2">
					Create new overlay for edit settings
				</NAlert>
			</template>
		</div>
	</div>
</template>

<style scoped>
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
</style>
