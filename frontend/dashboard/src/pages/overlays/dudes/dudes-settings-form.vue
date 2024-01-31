<script setup lang="ts">
import { Font, FontSelector } from '@twir/fontsource';
import { addZero, hexToRgb, colorBrightness } from '@zero-dependency/utils';
import { intervalToDuration } from 'date-fns';
import {
	NButton,
	NSlider,
	NColorPicker,
	useThemeVars,
	NDivider,
	NSwitch,
	NSelect,
	NDynamicTags,
	NTag,
	NInputNumber,
} from 'naive-ui';
import { h, computed, ref, watch } from 'vue';
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

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}

function formatDuration(duration: number) {
	const { hours = 0, minutes = 0, seconds = 0 } = intervalToDuration({ start: 0, end: duration });
	return `${addZero(hours)}:${addZero(minutes)}:${addZero(seconds)}`;
}

const fontData = ref<Font | null>(null);
watch(() => fontData.value, (font) => {
	if (!font) return;
	formValue.value.nameBoxSettings.fontFamily = font.id;
	formValue.value.nameBoxSettings.fontWeight = font.weights[0];
	formValue.value.messageBoxSettings.fontFamily = font.id;
});

const fontWeightOptions = computed(() => {
	if (!fontData.value) return [];
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }));
});
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
					<span>{{ t('overlays.dudes.dudeDefaultColor') }}</span>
					<n-color-picker
						v-model:value="formValue.dudeSettings.color"
						:modes="['hex']"
					/>
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

				<n-divider/>

				<div class="switch">
					<span>{{ t('overlays.dudes.dudeSounds') }}</span>
					<n-switch v-model:value="formValue.dudeSettings.soundsEnabled"/>
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

				<n-divider/>

				<div>
					<span>{{ t('overlays.dudes.nameBoxFill') }}</span>
					<n-dynamic-tags
						v-model:value="formValue.nameBoxSettings.fill"
						:max="3"
						:render-tag="(tag: string, index: number) => {
							const rgb = hexToRgb(tag)
							const textColor = rgb && colorBrightness(rgb) > 128 ? '#000' : '#fff'

							return h(NTag, {
								closable: true,
								onClose: () => {
									formValue.nameBoxSettings.fill.splice(index, 1)
								},
								style: {
									'--n-close-icon-color': textColor,
									'--n-close-icon-color-hover': textColor,
								},
								color: {
									color: tag,
									borderColor: tag,
									textColor,
								}
							}, { default: () => tag })
						}"
					>
						<template #input="{ submit, deactivate }">
							<n-color-picker
								:style="{
									position: 'absolute',
									width: '80px'
								}"
								size="small"
								default-show
								:show-alpha="false"
								:modes="['hex']"
								:actions="['confirm']"
								@confirm="submit($event)"
								@update-show="deactivate"
								@blur="deactivate"
							/>
						</template>
					</n-dynamic-tags>
				</div>

				<div>
					<span>{{ t('overlays.dudes.nameFillGradientStops') }}</span>
					<n-dynamic-tags
						v-model:value="formValue.nameBoxSettings.fillGradientStops"
						:on-create="(label) => {
							return Number(label)
						}"
						:render-tag="(tag: string, index: number) => {
							return h(NTag, {
								closable: true,
								onClose: () => {
									formValue.nameBoxSettings.fillGradientStops.splice(index, 1)
								}
							}, { default: () => tag })
						}"
						:max="formValue.nameBoxSettings.fill.length"
					>
						<template #input="{ submit, deactivate }">
							<n-input-number
								:style="{
									position: 'absolute',
									width: '100px'
								}"
								autofocus
								placeholder=""
								:max="1"
								:min="0"
								:step="0.01"
								:default-value="0.1"
								size="small"
								:update-value-on-input="false"
								:parse="(v) => {
									const parsedNum = Number(v)
									return Number.isNaN(parsedNum) ? 0 : parsedNum
								}"
								@keyup.enter="submit($event.target.value)"
								@confirm="submit($event)"
								@blur="deactivate"
							/>
						</template>
					</n-dynamic-tags>
				</div>

				<div>
					<span>{{ t('overlays.dudes.nameBoxGradientType') }}</span>
					<n-select
						v-model:value="formValue.nameBoxSettings.fillGradientType"
						:disabled="formValue.nameBoxSettings.fill.length < 2"
						:options="[
							{
								label: 'Vertical',
								value: 0,
							},
							{
								label: 'Horizontal',
								value: 1,
							}
						]"
					/>
				</div>

				<div>
					<span>{{ t('overlays.dudes.nameBoxFontFamily') }}</span>
					<font-selector
						v-model:selected-font="formValue.nameBoxSettings.fontFamily"
						:font-family="formValue.nameBoxSettings.fontFamily"
						:font-weight="formValue.nameBoxSettings.fontWeight"
						:font-style="formValue.nameBoxSettings.fontStyle"
						:subsets="['latin', 'cyrillic']"
						@update-font="(v) => fontData = v"
					/>
				</div>

				<div>
					<span>{{ t('overlays.dudes.nameBoxFontSize') }}</span>
					<n-slider
						v-model:value="formValue.nameBoxSettings.fontSize"
						:min="1"
						:max="128"
					/>
				</div>

				<div>
					<span>{{ t('overlays.dudes.nameBoxFontWeight') }}</span>
					<n-select
						v-model:value="formValue.nameBoxSettings.fontWeight"
						:options="fontWeightOptions"
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
