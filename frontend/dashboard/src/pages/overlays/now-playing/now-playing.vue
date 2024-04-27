<script setup lang="ts">
import { NowPlaying } from '@twir/frontend-now-playing';
import { NAlert, NResult, NTabPane, NTabs, useThemeVars, NA } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	useLastfmIntegration,
	useNowPlayingOverlayManager, useSpotifyIntegration,
	useUserAccessFlagChecker,
} from '@/api';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import { ChannelRolePermissionEnum } from '@/gql/graphql';
import NowPlayingForm from '@/pages/overlays/now-playing/now-playing-form.vue';
import {
	useNowPlayingForm,
	defaultSettings,
} from '@/pages/overlays/now-playing/use-now-playing-form';

const themeVars = useThemeVars();
const { t } = useI18n();
const { dialog } = useNaiveDiscrete();

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays);
const nowPlayingOverlayManager = useNowPlayingOverlayManager();
const creator = nowPlayingOverlayManager.useCreate();
const deleter = nowPlayingOverlayManager.useDelete();

const { data: spotifyData } = useSpotifyIntegration().useData();
const { data: lastFmData } = useLastfmIntegration().useData();

const isSomeSongIntegrationEnabled = computed(() => {
	return spotifyData.value?.userName || lastFmData.value?.userName;
});

const formStore = useNowPlayingForm();
const { data: settings } = storeToRefs(formStore);

const {
	data: entities,
} = nowPlayingOverlayManager.useGetAll();

const openedTab = ref<string>();

function resetTab() {
	if (!entities.value?.settings.at(0)) {
		openedTab.value = undefined;
		return;
	}

	openedTab.value = entities.value.settings.at(0)?.id;
}

async function handleAdd() {
	await creator.mutateAsync(defaultSettings);
}

async function handleClose(id: string) {
	dialog.create({
		title: 'Delete preset',
		content: 'Are you sure you want to delete this preset?',
		positiveText: 'Delete',
		negativeText: 'Cancel',
		showIcon: false,
		onPositiveClick: async () => {
			const entity = entities.value?.settings.find(s => s.id === id);
			if (!entity?.id) return;

			await deleter.mutateAsync(entity.id);
			resetTab();
		},
	});
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.settings.length ?? 0) < 5;
});

watch(openedTab, async (v) => {
	const entity = entities.value?.settings.find(s => s.id === v);
	if (!entity) return;
	formStore.$setData(entity);
});

watch(entities, () => {
	resetTab();
}, { immediate: true });
</script>

<template>
	<div class="flex flex-col gap-3">
		<div>
			<NowPlaying
				:settings="settings ?? { preset: 'TRANSPARENT' }"
				:track="{
					image_url: 'https://i.scdn.co/image/ab67616d0000b273e7fbc0883149094912559f2c',
					artist: 'Slipknot',
					title: 'Psychosocial'
				}"
			/>
		</div>
		<div>
			<n-result
				v-if="!isSomeSongIntegrationEnabled"
				status="warning"
				title="No enabled song integrations!"
			>
				<template #footer>
					Connect Spotify or Last.fm in
					<router-link :to="{ name: 'Integrations' }" #="{ navigate, href }" custom>
						<n-a :href="href" @click="navigate">
							{{ t('sidebar.integrations') }}
						</n-a>
					</router-link>
					to use this overlay
				</template>
			</n-result>
			<template v-if="isSomeSongIntegrationEnabled">
				<n-tabs
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
					<template v-if="entities?.settings.length">
						<n-tab-pane
							v-for="(entity, entityIndex) in entities?.settings"
							:key="entity.id"
							:tab="`#${entityIndex+1}`"
							:name="entity.id!"
						>
							<now-playing-form />
						</n-tab-pane>
					</template>
				</n-tabs>
				<n-alert v-if="!entities?.settings.length" type="info" class="mt-2">
					Create new overlay for edit settings
				</n-alert>
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
