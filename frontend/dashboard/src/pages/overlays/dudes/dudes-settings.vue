<script setup lang="ts">
import {
	NTabs,
	NTabPane,
	NAlert,
	useThemeVars,
	NButton,
	NSpace,
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
import CommandButton from '@/components/commandButton.vue';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';

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
	if (!profile.value || !openedTab.value) return null;
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
	<div style="display: flex; gap: 42px; height: calc(100% - var(--layout-header-height));">
		<div style="width: 70%; position: relative;">
			<div v-if="dudesIframeUrl">
				<iframe
					ref="dudesIframe"
					style="position: absolute;"
					:src="dudesIframeUrl"
					class="iframe"
				/>
				<n-space :size="6" style="position: absolute; top: 18px; left: 8px;">
					<n-button @click="dudesIframeStore.sendIframeMessage('spawn-emote')">
						Emote
					</n-button>
					<n-button @click="dudesIframeStore.sendIframeMessage('show-message')">
						Message
					</n-button>
					<n-button @click="dudesIframeStore.sendIframeMessage('jump')">
						Jump
					</n-button>
					<n-button @click="dudesIframeStore.sendIframeMessage('grow')">
						Grow
					</n-button>
					<n-button @click="dudesIframeStore.sendIframeMessage('leave')">
						Leave
					</n-button>
					<n-button @click="dudesIframeStore.sendIframeMessage('reset')">
						Reset
					</n-button>
				</n-space>
			</div>
		</div>
		<div style="width: 30%;">
			<div style="display: flex; gap: 8px; flex-wrap: wrap">
				<command-button title="Jump command" name="jump" />
				<command-button title="Color command" name="dudes color" />
				<command-button title="Grow command" name="dudes grow" />
				<command-button title="Sprite command" name="dudes sprite" />
				<command-button title="Leave command" name="dudes leave" />
			</div>
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
						<dudes-settings-form />
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
@import '../styles.css';

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
