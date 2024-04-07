<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { TwirEventType } from '@twir/api/messages/events/events';
import type {
	Settings_AnimationSettings,
} from '@twir/api/messages/overlays_kappagen/overlays_kappagen';
import { useDebounceFn } from '@vueuse/core';
import { CopyIcon } from 'lucide-vue-next';
import { NTabs, NTabPane, NButton, NButtonGroup, useThemeVars } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import SettingsAnimations from './settingsAnimations.vue';
import SettingsEvents from './settingsEvents.vue';
import SettingsGeneral from './settingsGeneral.vue';
import { useKappagenFormSettings } from './store.js';


import { useKappaGenOverlayManager, useProfile } from '@/api';
import { flatEvents } from '@/components/events/helpers.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';
import { useToast } from '@/components/ui/toast';

const availableEvents = Object.values(flatEvents)
	.filter(e => e.enumValue !== undefined && TwirEventType[e.enumValue])
	.map(e => {
		return {
			name: e.name,
			value: e.enumValue,
		};
	}) as Array<{ name: string, value: TwirEventType }>;

const themeVars = useThemeVars();
const { t } = useI18n();
const { copyOverlayLink } = useCopyOverlayLink('kappagen');

const kappagenManager = useKappaGenOverlayManager();
const { data: settings, error } = kappagenManager.getSettings();
const updater = kappagenManager.updateSettings();

const { settings: formValue } = useKappagenFormSettings();

watch(error, async (v) => {
	if (v instanceof RpcError) {
		if (v.code === 'not_found') {
			await updater.mutateAsync(formValue.value);
		}
	}
});

const { toast } = useToast();
async function save() {
	if (!formValue.value) return;

	await updater.mutateAsync(formValue.value);
	toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
	});
}

const sendSettings = () => sendIframeMessage('settings', {
	...toRaw(formValue.value),
	channelName: profile.value?.login,
	channelId: profile.value?.id,
});

const debouncedSave = useDebounceFn(async () => {
	await save();
	sendSettings();
}, 1000);
watch(formValue, () => {
	debouncedSave();
}, { deep: true });

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

watch(kappagenIframeRef, (v) => {
	if (!v) return;
	v.contentWindow?.addEventListener('message', (event) => {
		const data = JSON.parse(event.data);
		if (data.key !== 'getSettings') return;
		sendSettings();
	});
});

const playKappaPreview = (animation: Settings_AnimationSettings) => {
	sendIframeMessage('kappaWithAnimation', { animation });
};
</script>

<template>
	<div class="flex h-full p-6 gap-10">
		<div class="w-1/2">
			<div class="header-buttons">
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
						<CopyIcon class="mr-2 h-6 w-6" />
						{{ t('overlays.copyOverlayLink') }}
					</n-button>
				</n-button-group>
			</div>

			<n-tabs
				default-value="main"
				type="line"
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

		<div class="w-1/2 h-full">
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
.header-buttons {
	display: flex;
	justify-content: space-between;
	flex-wrap: wrap-reverse;
	row-gap: 10px;
	column-gap: 10px;
}

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
