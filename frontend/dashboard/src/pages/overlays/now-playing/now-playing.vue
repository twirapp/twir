<script setup lang="ts">
import { NowPlaying } from '@twir/frontend-now-playing';
import { NAlert, NTabPane, NTabs, useThemeVars } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	useNowPlayingOverlayManager,
	useUserAccessFlagChecker,
} from '@/api';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import NowPlayingForm from '@/pages/overlays/now-playing/now-playing-form.vue';
import { useNowPlayingForm } from '@/pages/overlays/now-playing/use-now-playing-form';


const themeVars = useThemeVars();
const { t } = useI18n();
const { dialog } = useNaiveDiscrete();

const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const nowPlayingOverlayManager = useNowPlayingOverlayManager();
const creator = nowPlayingOverlayManager.useCreate();
const deleter = nowPlayingOverlayManager.useDelete();

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
	await creator.mutateAsync();
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
	<div style="display: flex; gap: 42px; height: calc(100% - var(--layout-header-height));">
		<div style="width: 70%; position: relative;">
			<div v-if="settings" class="iframe" style="padding: 10px">
				<NowPlaying
					:settings="settings"
					:track="{artist: '123', title: 'asd'}"
				/>
			</div>
		</div>
		<div style="width: 30%;">
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
			<n-alert v-if="!entities?.settings.length" type="info" style="margin-top: 8px;">
				Create new overlay for edit settings
			</n-alert>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

.iframe {
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
}
</style>
