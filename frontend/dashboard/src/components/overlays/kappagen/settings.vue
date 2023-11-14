<script lang="ts" setup>
import { TwirEventType } from '@twir/grpc/generated/api/api/events';
import type {
	Settings, Settings_AnimationSettings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification, NTabs, NTabPane, NButton, NButtonGroup, NSlider, NSwitch, NDivider, NCheckboxGroup, NCheckbox, NGrid, NGridItem, NAlert } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import animationSettings from './animationSettings.vue';
import { animations } from './kappagen_animations';

import { useKappaGenOverlayManager, useProfile } from '@/api';
import CommandButton from '@/components/commandButton.vue';
import { flatEvents } from '@/components/events/helpers.js';

const availableEvents = Object.values(flatEvents)
	.filter(e => e.enumValue !== undefined && TwirEventType[e.enumValue])
	.map(e => {
		return {
			name: e.name,
			value: e.enumValue,
		};
	}) as Array<{ name: string, value: TwirEventType }>;

const formValue = ref<Settings>({
	emotes: {
		time: 5,
		max: 0,
		queue: 0,
	},
	animations: animations,
	enableRave: false,
	animation: {
		fadeIn: true,
		fadeOut: true,
		zoomIn: true,
		zoomOut: true,
	},
	cube: {
		speed: 6,
	},
	size: {
		// from 7 to 20
		ratioNormal: 7,
		// from 14 to 40
		ratioSmall: 14,
		min: 1,
		max: 256,
	},
	enabledEvents: [],
	enableSpawn: true,
});
const kappagenManager = useKappaGenOverlayManager();
const { data: settings } = kappagenManager.getSettings();
watch(settings, (s) => {
	if (!s) return;

	formValue.value = toRaw(s);
});

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
	<div style="display: flex; gap: 15px;">
		<iframe
			v-if="kappagenIframeUrl"
			ref="kappagenIframeRef"
			:src="kappagenIframeUrl"
			style="width: 50%; height: auto; aspect-ratio: 16/9; border: 0;"
		/>

		<div style="width: 50%">
			<div style="display: flex; justify-content: space-between;">
				<n-button-group style="width: 100%">
					<n-button secondary @click="sendIframeMessage('kappa', 'EZ')">
						{{ t('overlays.kappagen.testKappagen') }}
					</n-button>
					<n-button secondary type="info" @click="sendIframeMessage('spawn', ['EZ'])">
						{{ t('overlays.kappagen.testSpawn') }}
					</n-button>
				</n-button-group>

				<n-button secondary type="success" @click="save">
					{{ t('sharedButtons.save') }}
				</n-button>
			</div>

			<n-alert title="Info" type="info" :show-icon="false" style="margin-top: 5px;">
				{{ t('overlays.kappagen.info') }}
			</n-alert>

			<n-tabs default-value="main" type="line" size="large" justify-content="space-evenly" animated style="width: 100%">
				<n-tab-pane name="main" :tab="t('overlays.kappagen.tabs.main')">
					<div class="tab">
						<CommandButton name="kappagen" />

						<div class="switch">
							<n-switch v-model:value="formValue.enableSpawn" />
							<span>{{ t('overlays.kappagen.settings.spawn') }}</span>
						</div>

						<n-divider />

						<div class="slider">
							{{ t('overlays.kappagen.settings.size') }}({{ formValue.size!.ratioNormal }})
							<n-slider
								v-model:value="formValue.size!.ratioNormal"
								reverse
								:min="7"
								:max="20"
							/>
						</div>

						<div class="slider">
							{{ t('overlays.kappagen.settings.sizeSmall') }}({{ formValue.size!.ratioSmall }})
							<n-slider
								v-model:value="formValue.size!.ratioSmall"
								reverse
								:min="14"
								:max="40"
							/>
						</div>

						<n-divider />

						<div class="slider">
							{{ t('overlays.kappagen.settings.time') }}({{ formValue.emotes!.time }}s)
							<n-slider
								v-model:value="formValue.emotes!.time"
								:min="1"
								:max="15"
							/>
						</div>

						<div class="slider">
							{{ t('overlays.kappagen.settings.maxEmotes') }}({{ formValue.emotes!.max }})
							<n-slider
								v-model:value="formValue.emotes!.max"
								:min="0"
								:max="250"
							/>
						</div>

						<n-divider />

						<div class="switchers">
							<span>{{ t('overlays.kappagen.settings.animationsOnAppear') }}</span>

							<div class="switch">
								<n-switch v-model:value="formValue.animation!.fadeIn" />
								<span>Fade</span>
							</div>

							<div class="switch">
								<n-switch v-model:value="formValue.animation!.zoomIn" />
								<span>Zoom</span>
							</div>
						</div>

						<n-divider />

						<div class="switchers">
							<span>{{ t('overlays.kappagen.settings.animationsOnDisappear') }}</span>

							<div class="switch">
								<n-switch v-model:value="formValue.animation!.fadeOut" />
								<span>Fade</span>
							</div>

							<div class="switch">
								<n-switch v-model:value="formValue.animation!.zoomOut" />
								<span>Zoom</span>
							</div>
						</div>

						<n-divider />

						<div class="switch">
							<n-switch v-model:value="formValue.enableRave" />
							<span>{{ t('overlays.kappagen.settings.rave') }}</span>
						</div>
					</div>
				</n-tab-pane>

				<n-tab-pane name="events" :tab="t('overlays.kappagen.tabs.events')">
					<n-checkbox-group v-model:value="formValue.enabledEvents">
						<div style="display: flex; flex-direction: column; gap: 5px;">
							<n-checkbox
								v-for="event of availableEvents"
								:key="event.name"
								:value="event.value"
								:label="event.name"
							/>
						</div>
					</n-checkbox-group>
				</n-tab-pane>

				<n-tab-pane name="animations" :tab="t('overlays.kappagen.tabs.animations')">
					<n-grid :cols="2" :x-gap="16" :y-gap="16" responsive="self">
						<n-grid-item v-for="animation of formValue.animations" :key="animation.style" :span="1">
							<animationSettings :settings="animation" @play="playKappaPreview" />
						</n-grid-item>
					</n-grid>
				</n-tab-pane>
			</n-tabs>
		</div>
	</div>
</template>

<style scoped>
.tab {
	display: flex;
	flex-direction: column;
	gap: 15px;
}

.slider {
	display: flex;
	gap: 5px;
	flex-direction: column;
}

.switchers {
	display: flex;
	gap: 5px;
	flex-direction: column;
}

.switch {
	display: flex;
	gap: 5px;
}

:deep(.n-divider) {
	margin-top: 0;
	margin-bottom: 0;
}
</style>
