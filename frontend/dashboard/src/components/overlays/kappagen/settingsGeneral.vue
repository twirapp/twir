<script lang="ts" setup>
import { NSlider, NSwitch, NAlert, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useSettings } from './store.js';

import CommandButton from '@/components/commandButton.vue';

const { settings: formValue } = useSettings();
const { t } = useI18n();

const formatSizeValue = (v: number) => parseInt(`${v}`.split('.')[1]);
</script>

<template>
	<div class="tab">
		<n-alert type="info" :show-icon="false" style="margin-top: 5px;">
			{{ t('overlays.kappagen.info') }}
		</n-alert>
		<CommandButton name="kappagen" />

		<div class="switch">
			<n-switch v-model:value="formValue.enableSpawn" />
			<span>{{ t('overlays.kappagen.settings.spawn') }}</span>
		</div>

		<n-divider />

		<div class="slider">
			{{ t('overlays.kappagen.settings.size') }}({{ formatSizeValue(formValue.size!.ratioNormal) }})
			<n-slider
				v-model:value="formValue.size!.ratioNormal"
				:format-tooltip="formatSizeValue"
				:step="0.01"
				:min="0.05"
				:max="0.15"
			/>
		</div>

		<div class="slider">
			{{ t('overlays.kappagen.settings.sizeSmall') }}({{ formatSizeValue(formValue.size!.ratioSmall) }})
			<n-slider
				v-model:value="formValue.size!.ratioSmall"
				:format-tooltip="formatSizeValue"
				:step="0.01"
				:min="0.02"
				:max="0.07"
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
</template>
