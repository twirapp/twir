<script setup lang="ts">
import {
	NTabs,
	NTabPane,
	NAlert,
	NScrollbar,
	useThemeVars,
} from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import DudesSettingsForm from './dudes-settings-form.vue';
import { DudesSettingsWithOptionalId } from './dudes-settings.js';
import { useDudesForm } from './use-dudes-form.js';
import { useDudesIframe } from './use-dudes-frame.js';

import {
	useDudesOverlayManager, useProfile, useUserAccessFlagChecker,
} from '@/api/index.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';
import CommandButton from '@/components/commandButton.vue';

const themeVars = useThemeVars();
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const dudesOverlayManager = useDudesOverlayManager();
const creator = dudesOverlayManager.useCreate();
const deleter = dudesOverlayManager.useDelete();
const { data: profile } = useProfile();

const { t } = useI18n();
const { dialog } = useNaiveDiscrete();

const {
	data: entities,
} = dudesOverlayManager.useGetAll();

const openedTab = ref<string>();

const dudesIframeStore = useDudesIframe();
const { dudesIframe } = storeToRefs(dudesIframeStore);
const dudesIframeUrl = computed(() => {
	if (!profile.value) return null;
	return `${window.location.origin}/overlays/${profile.value.apiKey}/dudes?id=${openedTab.value}`;
});

const { $setData, $getDefaultSettings } = useDudesForm();

function resetTab() {
	if (!entities.value?.settings.at(0)) {
		openedTab.value = undefined;
		return;
	}

	openedTab.value = entities.value.settings.at(0)?.id;
}

watch(entities, () => {
	resetTab();
}, { immediate: true });

watch(openedTab, (v) => {
	const entity = entities.value?.settings.find(s => s.id === v) as DudesSettingsWithOptionalId;
	if (!entity) return;
	$setData(entity);
});

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

async function handleAdd() {
	await creator.mutateAsync($getDefaultSettings());
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.settings.length ?? 0) < 5;
});
</script>

<template>
	<div style="display: flex; gap: 42px; height: 100%; padding: 24px;">
		<div style="width: 50%">
			<iframe
				v-if="dudesIframeUrl"
				ref="dudesIframe"
				:src="dudesIframeUrl"
				class="iframe"
			/>
		</div>
		<div style="width: 50%; height: 100%;">
			<command-button name="jump"/>
			<n-tabs
				v-model:value="openedTab"
				type="card"
				:closable="userCanEditOverlays"
				:addable="addable"
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
						<n-scrollbar style="max-height: 75vh;" trigger="none">
							<dudes-settings-form/>
						</n-scrollbar>
					</n-tab-pane>
				</template>
			</n-tabs>
			<n-alert v-if="!entities?.settings.length" type="info" style="margin-top: 8px;">
				Create new overlay for edit settings
			</n-alert>
		</div>
	</div>
</template>

<style scope>
.iframe {
	height: 100%;
	width: 100%;
	aspect-ratio: 16 / 9;
	border: 0;
	margin-top: 8px;
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
}
</style>
