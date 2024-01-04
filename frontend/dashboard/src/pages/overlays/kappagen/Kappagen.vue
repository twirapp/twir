<script lang="ts" setup>
import { TwirEventType } from '@twir/grpc/generated/api/api/events';
import type {
	Settings_AnimationSettings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification, NTabs, NTabPane, NButton, NButtonGroup, useThemeVars } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import SettingsAnimations from './settingsAnimations.vue';
import SettingsEvents from './settingsEvents.vue';
import SettingsGeneral from './settingsGeneral.vue';
import { useSettings } from './store.js';


import { useKappaGenOverlayManager, useProfile } from '@/api';
import { flatEvents } from '@/components/events/helpers.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';

const availableEvents = Object.values(flatEvents)
	.filter(e => e.enumValue !== undefined && TwirEventType[e.enumValue])
	.map(e => {
		return {
			name: e.name,
			value: e.enumValue,
		};
	}) as Array<{ name: string, value: TwirEventType }>;

const themeVars = useThemeVars();

const { copyOverlayLink } = useCopyOverlayLink('kappagen');

const { settings: formValue } = useSettings();

const kappagenManager = useKappaGenOverlayManager();
const { data: settings } = kappagenManager.getSettings();
watch(settings, (s) => {
	if (!s) return;

	const events = toRaw(s.events);

	for (const event of availableEvents) {
		const isExists = events.some(e => e.event === event.value);
		if (isExists) continue;

		events.push({ event: event.value, disabledStyles: [], enabled: false });
	}

	formValue.value = {
		...toRaw(s),
		events: events,
	};
}, { immediate: true });

watch(() => [
	formValue.value.emotes,
	formValue.value.enableRave,
	formValue.value.animation,
	formValue.value.cube,
	formValue.value.size,
], () => sendSettings(), { deep: true });

const updater = kappagenManager.updateSettings();
const { data: profile } = useProfile();

const kappagenIframeRef = ref<HTMLIFrameElement | null>(null);
const kappagenIframeUrl = computed(() => {
	if (!profile.value) return null;

	return `${window.location.origin}/overlays/${profile.value.apiKey}/kappagen`;
});

const sendIframeMessage = (key: string, data?: any) => {
	if (!kappagenIframeRef.value) return;
	const win = kappagenIframeRef.value;

	win.contentWindow?.postMessage(JSON.stringify({
		key,
		data: toRaw(data),
	}));
};

const sendSettings = () => sendIframeMessage('settings', {
	...toRaw(formValue.value),
	channelName: profile.value?.login,
	channelId: profile.value?.id,
});

watch(kappagenIframeRef, (v) => {
	if (!v) return;
	v.contentWindow?.addEventListener('message', (e) => {
		if (e.data !== 'getSettings') return;

		sendSettings();
	});
});

const message = useNotification();
const { t } = useI18n();

const playKappaPreview = (animation: Settings_AnimationSettings) => {
	sendIframeMessage('kappaWithAnimation', { animation });
};

async function save() {
	if (!formValue.value) return;

	await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}
</script>

<template>
	<div style="display: flex; gap: 42px; height: 100%; padding: 24px;">
		<div style="width: 50%">
			<div style="display: flex; justify-content: space-between;">
				<n-button-group>
					<n-button secondary @click="sendIframeMessage('kappa', 'EZ')">
						{{ t('overlays.kappagen.testKappagen') }}
					</n-button>
					<n-button secondary type="info" @click="sendIframeMessage('spawn', ['EZ'])">
						{{ t('overlays.kappagen.testSpawn') }}
					</n-button>

					<n-button secondary type="warning" @click="sendIframeMessage('clear')">
						{{ t('overlays.kappagen.clear') }}
					</n-button>
				</n-button-group>

				<n-button-group>
					<n-button secondary type="info" @click="copyOverlayLink()">
						{{ t('overlays.copyOverlayLink') }}
					</n-button>
					<n-button secondary type="success" @click="save">
						{{ t('sharedButtons.saveSettings') }}
					</n-button>
				</n-button-group>
			</div>

			<n-tabs
				default-value="main"
				type="segment"
				size="large"
				justify-content="space-evenly"
				animated
				style="width: 100%; margin-top: 16px;"
			>
				<n-tab-pane name="main" :tab="t('overlays.kappagen.tabs.main')">
					<div class="card">
						<div class="content">
							<SettingsGeneral />
						</div>
					</div>
				</n-tab-pane>

				<n-tab-pane name="events" :tab="t('overlays.kappagen.tabs.events')">
					<div class="card">
						<div class="content">
							<SettingsEvents />
						</div>
					</div>
				</n-tab-pane>

				<n-tab-pane name="animations" :tab="t('overlays.kappagen.tabs.animations')">
					<div class="card">
						<div class="content">
							<SettingsAnimations @play="playKappaPreview" />
						</div>
					</div>
				</n-tab-pane>
			</n-tabs>
		</div>

		<div style="width: 50%; height: 100%;">
			<iframe
				v-if="kappagenIframeUrl"
				ref="kappagenIframeRef"
				:src="kappagenIframeUrl"
				class="iframe"
			/>
		</div>
	</div>
</template>

<style scoped>
:deep(.card) {
	display: flex;
	flex-direction: column;
	gap: 8px;
	height: 100%;
	border-radius: 4px;
	background-color: v-bind('themeVars.actionColor');
}

:deep(.card .content) {
	padding: 12px;
}

:deep(.card .content .settings) {
	padding-top: 5px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

:deep(.card .title) {
	display: flex;
	justify-content: space-between;
	width: 100%;
	padding-bottom: 3px;
}

:deep(.card .title .info) {
	display: flex;
	gap: 4px;
}

:deep(.card .title-bordered) {
	border-bottom: 1px solid v-bind('themeVars.borderColor');
}

:deep(.card .form-item) {
	display: flex;
	justify-content: space-between;
	gap: 4px;
}

:deep(.n-input-number) {
	width: 40%
}

:deep(.tab) {
	display: flex;
	flex-direction: column;
	gap: 15px;
}

:deep(.slider) {
	display: flex;
	gap: 5px;
	flex-direction: column;
}

:deep(.switchers) {
	display: flex;
	gap: 5px;
	flex-direction: column;
}

:deep(.switch) {
	display: flex;
	gap: 5px;
}

:deep(.n-divider) {
	margin-top: 0;
	margin-bottom: 0;
}

:deep(.card) {
	background-color: v-bind('themeVars.actionColor');
}

.iframe {
	height: 100%;
	width: 100%;
	aspect-ratio: 16/9;
	border: 0;
	margin-top: 8px;
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
}
</style>
