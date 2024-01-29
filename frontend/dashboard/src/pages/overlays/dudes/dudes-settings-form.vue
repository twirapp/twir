<script setup lang="ts">
import { intervalToDuration } from 'date-fns';
import {
	NButton,
	NSlider,
	NColorPicker,
	useThemeVars,
	NDivider,
	NSwitch,
} from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useDudesForm } from './use-dudes-form.js';
import { useDudesIframe } from './use-dudes-frame.js';

import { useDudesOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';

const { t } = useI18n();
const themeVars = useThemeVars();
const discrete = useNaiveDiscrete();
const { copyOverlayLink } = useCopyOverlayLink('dudes');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = useProfile();

const { data: formValue, $reset } = useDudesForm();

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});

const dudesIframeStore = useDudesIframe();
const manager = useDudesOverlayManager();
const updater = manager.useUpdate();

async function save() {
	if (!formValue.value.id) return;

	await updater.mutateAsync({
		id: formValue.value.id,
		settings: formValue.value,
	});

	dudesIframeStore.dudesIframe?.contentWindow?.location.reload();

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}

function zeroPad(num: number) {
	return String(num).padStart(2, '0');
}

function formatDuration(duration: number) {
	const { hours = 0,  minutes = 0, seconds = 0 } = intervalToDuration({ start: 0, end: duration });
	return `${zeroPad(hours)}:${zeroPad(minutes)}:${zeroPad(seconds)}`;
}
</script>

<template>
	<div v-if="formValue" class="card">
		<div class="card-header">
			<n-button
				secondary
				type="error"
				@click="$reset"
			>
				{{ t('sharedButtons.setDefaultSettings') }}
			</n-button>
			<n-button
				secondary
				type="info"
				:disabled="!formValue.id || !canCopyLink"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</n-button>
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>

		<div class="card-body">
			<div class="card-body-column">
				<div style="margin-top: 1rem;">
					<span>{{ t('overlays.dudes.dudeColor') }}</span>
					<n-color-picker v-model:value="formValue.dudeSettings.color" />
				</div>

				<div>
					<span>{{ t('overlays.dudes.dudeGravity') }}</span>
					<n-slider
						v-model:value="formValue.dudeSettings.gravity"
						:min="100"
						:max="5000"
					/>
				</div>

				<div>
					<span>{{ t('overlays.dudes.dudeMaxLifeTime') }}</span>
					<n-slider
						v-model:value="formValue.dudeSettings.maxLifeTime"
						:min="1000"
						:max="120 * 60 * 1000"
						:step="1000"
						:format-tooltip="(value) => formatDuration(value)"
					/>
				</div>

				<div>
					<span>{{ t('overlays.dudes.dudeScale') }}</span>
					<n-slider
						v-model:value="formValue.dudeSettings.scale"
						:min="1"
						:max="10"
						:step="1"
					/>
				</div>

				<n-divider />

				<div class="switch">
					<span>{{ t('overlays.dudes.dudeSounds') }}</span>
					<n-switch v-model:value="formValue.dudeSettings.soundsEnabled" />
				</div>

				<div>
					<span>{{ t('overlays.dudes.dudeSoundsVolume') }}</span>
					<n-slider
						v-model:value="formValue.dudeSettings.soundsVolume"
						:min="0.01"
						:max="1"
						:step="0.01"
						:format-tooltip="(value) => `${(value * 100).toFixed(0)}%`"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

.card-header {
	flex-wrap: wrap;
	justify-content: flex-start;
}

.card-body-column {
	width: 100%;
	padding-bottom: 1rem;
}

.switch {
	display: flex;
	justify-content: space-between;
}

.card {
	background-color: v-bind('themeVars.cardColor');
}
</style>
