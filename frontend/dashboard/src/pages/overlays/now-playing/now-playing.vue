<script setup lang="ts">
import { NAlert, NTabPane, NTabs } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	useNowPlayingOverlayManager,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import NowPlayingForm from '@/pages/overlays/now-playing/now-playing-form.vue';
import { useNowPlayingIframe } from '@/pages/overlays/now-playing/now-playing-iframe';
import { useNowPlayingForm } from '@/pages/overlays/now-playing/use-now-playing-form';

const { t } = useI18n();
const { dialog } = useNaiveDiscrete();

const { data: profile } = useProfile();
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const nowPlayingOverlayManager = useNowPlayingOverlayManager();
const creator = nowPlayingOverlayManager.useCreate();
const deleter = nowPlayingOverlayManager.useDelete();

const iframeStore = useNowPlayingIframe();
const { nowPlayingIframe } = storeToRefs(iframeStore);
const { $setData } = useNowPlayingForm();

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

const overlayIframeUrl = computed(() => {
	if (!profile.value || !openedTab.value) return null;
	return `${window.location.origin}/overlays/${profile.value.apiKey}/now-playing?id=${openedTab.value}`;
});

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
	$setData(entity);
});

watch(entities, () => {
	resetTab();
}, { immediate: true });
</script>

<template>
	<div style="display: flex; gap: 42px; height: calc(100% - var(--layout-header-height));">
		<div style="width: 70%; position: relative;">
			<div v-if="overlayIframeUrl">
				<iframe
					ref="nowPlayingIframe"
					style="position: absolute;"
					:src="overlayIframeUrl"
					class="iframe"
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
</style>
