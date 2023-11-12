<script lang="ts" setup>
import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification, NTabs, NTabPane, NButton, NButtonGroup, NSlider, NSwitch, NDivider } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { animations } from './kappagen_animations';

import { useKappaGenOverlayManager, useProfile } from '@/api';

const formValue = ref<Required<Settings>>({
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
});
const kappagenManager = useKappaGenOverlayManager();
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

watch(formValue, (v) => {
	sendIframeMessage('settings', v);
}, { deep: true });

const message = useNotification();
const { t } = useI18n();

async function save() {
	if (!formValue.value) return;

	await updater.mutateAsync(formValue.value as Settings);
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
						Test kappagen
					</n-button>
					<n-button secondary type="info" @click="sendIframeMessage('spawn', ['EZ'])">
						Test spawn
					</n-button>
				</n-button-group>

				<n-button secondary type="success" @click="save">
					Save settings
				</n-button>
			</div>

			<n-tabs default-value="main" type="line" size="large" justify-content="space-evenly" animated style="width: 100%">
				<n-tab-pane name="main" tab="Main settings">
					<div class="tab">
						<div class="slider">
							Size {{ formValue.size.ratioNormal }}
							<n-slider
								v-model:value="formValue.size.ratioNormal"
								reverse
								:min="7"
								:max="20"
							/>
						</div>

						<div class="slider">
							Small size {{ formValue.size.ratioSmall }}
							<n-slider
								v-model:value="formValue.size.ratioSmall"
								reverse
								:min="14"
								:max="40"
							/>
						</div>

						<div class="slider">
							The time an emote stays on screen, in seconds ({{ formValue.emotes.time }})
							<n-slider
								v-model:value="formValue.emotes.time"
								:min="1"
								:max="15"
							/>
						</div>

						<n-divider />

						<div class="switchers">
							<span>Show animations</span>

							<div class="switch">
								<n-switch v-model:value="formValue.animation.fadeIn" />
								<span>Fade</span>
							</div>

							<div class="switch">
								<n-switch v-model:value="formValue.animation.zoomIn" />
								<span>Zoom</span>
							</div>
						</div>

						<n-divider />

						<div class="switchers">
							<span>Hide animations</span>

							<div class="switch">
								<n-switch v-model:value="formValue.animation.fadeOut" />
								<span>Fade</span>
							</div>

							<div class="switch">
								<n-switch v-model:value="formValue.animation.zoomOut" />
								<span>Zoom</span>
							</div>
						</div>

						<n-divider />

						<div class="switch">
							<n-switch v-model:value="formValue.enableRave" />
							<span>Rave</span>
						</div>
					</div>
				</n-tab-pane>

				<n-tab-pane name="animations" tab="Animations">
					{{ formValue.animations }}
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
