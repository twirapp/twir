<script setup lang="ts">
import { Font, FontSelector } from '@twir/fontsource';
import { DudesSprite } from '@twir/types/overlays';
import { addZero, hexToRgb, colorBrightness, capitalize } from '@zero-dependency/utils';
import { intervalToDuration } from 'date-fns';
import {
	NButton,
	NSlider,
	NColorPicker,
	useThemeVars,
	NSwitch,
	NSelect,
	NDynamicTags,
	NTag,
	NInputNumber,
	NFormItem,
	NScrollbar,
	NForm,
	NTabs,
	NTabPane,
} from 'naive-ui';
import { h, computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useDudesForm } from './use-dudes-form.js';
import { useDudesIframe } from './use-dudes-frame.js';

import { useDudesOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api/index.js';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';
import SelectTwitchUsers from '@/components/twitchUsers/multiple.vue';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';


const { t } = useI18n();
const themeVars = useThemeVars();
const discrete = useNaiveDiscrete();
const { copyOverlayLink } = useCopyOverlayLink('dudes');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = useProfile();

const { data: formValue, $reset } = useDudesForm();
const { sendIframeMessage } = useDudesIframe();

watch(formValue, (form) => {
	if (!form) return;
	if (!form.nameBoxSettings.fill.length) return;
	sendIframeMessage('update-settings', form);
}, { deep: true });

watch(() => formValue.value.dudeSettings.defaultSprite, (dudeSprite) => {
	sendIframeMessage('update-sprite', dudeSprite);
});

watch(() => formValue.value.dudeSettings.color, (dudeColor) => {
	sendIframeMessage('update-color', dudeColor);
});

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});

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
	if (hours === 0) {
		return `${addZero(minutes)}:${addZero(seconds)}`;
	}

	return `${addZero(hours)}:${addZero(minutes)}:${addZero(seconds)}`;
}

const fillGradientStops = computed(() => {
	if (!formValue.value) return [];
	return formValue.value.nameBoxSettings.fillGradientStops.map((stop) => `${stop}`);
});

const fillGradidentStopMessage = computed(() => {
	if (!formValue.value.nameBoxSettings.fillGradientStops.length) {
		return t('overlays.dudes.nameBoxFillGradientStopsError');
	}

	return '';
});

const nameBoxFillMessage = computed(() => {
	if (!formValue.value.nameBoxSettings.fill.length) {
		return t('overlays.dudes.nameBoxFillError');
	}

	return '';
});

const fontData = ref<Font | null>(null);
watch(() => fontData.value, (font) => {
	if (!font) return;
	formValue.value.nameBoxSettings.fontFamily = font.id;
	formValue.value.messageBoxSettings.fontFamily = font.id;
});

const fontWeightOptions = computed(() => {
	if (!fontData.value) return [];
	return fontData.value.weights.map((weight) => ({
		label: `${weight}`,
		value: weight,
	}));
});

const fontStyleOptions = computed(() => {
	if (!fontData.value) return [];
	return fontData.value.styles.map((style) => ({
		label: capitalize(style),
		value: style,
	}));
});

const fontVariantOptions = ['normal', 'small-caps'].map((variant) => ({
	label: capitalize(variant),
	value: variant,
}));

const lineJoinOptions = ['round', 'bevel', 'miter'].map((lineJoin) => ({
	label: capitalize(lineJoin),
	value: lineJoin,
}));

const isMessageBoxDisabled = computed(() => {
	return !formValue.value.messageBoxSettings.enabled;
});

const isNameBoxDisabled = computed(() => {
	return !formValue.value.dudeSettings.visibleName;
});

const isDropShadowDisabled = computed(() => {
	return isNameBoxDisabled.value || !formValue.value.nameBoxSettings.dropShadow;
});

const dudesSprites = Object.keys(DudesSprite).map((key) => ({
	label: capitalize(key),
	value: key,
}));
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

			<n-tabs class="pt-2" type="line" placement="left" default-value="dude" animated>
				<n-tab-pane name="dude" :tab="t('overlays.dudes.dudeDivider')">
					<n-form-item :label="t('overlays.dudes.dudeDefaultSprite')">
						<n-select
							v-model:value="formValue.dudeSettings.defaultSprite"
							:options="dudesSprites"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.dudeMaxOnScreen')">
						<n-slider
							v-model:value="formValue.dudeSettings.maxOnScreen"
							:min="0"
							:max="128"
							:step="1"
							:format-tooltip="(value) => {
								if (value === 0) {
									return t('overlays.dudes.dudeMaxOnScreenUnlimited');
								}

								return value;
							}"
						/>
					</n-form-item>

					<n-form-item :label="t('overlays.dudes.dudeColor')">
						<n-color-picker
							v-model:value="formValue.dudeSettings.color"
							:modes="['hex']"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.dudeGravity')">
						<n-slider
							v-model:value="formValue.dudeSettings.gravity"
							:min="100"
							:max="5000"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.dudeMaxLifeTime')">
						<n-slider
							v-model:value="formValue.dudeSettings.maxLifeTime"
							:min="1000"
							:max="120 * 60 * 1000"
							:step="1000"
							:format-tooltip="(value) => formatDuration(value)"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.dudeScale')">
						<n-slider
							v-model:value="formValue.dudeSettings.scale"
							:min="1"
							:max="10"
							:step="1"
						/>
					</n-form-item>
				</n-tab-pane>

				<n-tab-pane name="ignoring" :tab="t('overlays.dudes.ignoreDivider')">
					<n-form-item
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.ignoreCommands')"
					>
						<n-switch v-model:value="formValue.ignoreSettings.ignoreCommands" />
					</n-form-item>

					<n-form-item
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.ignoreUsers')"
					>
						<n-switch v-model:value="formValue.ignoreSettings.ignoreUsers" />
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.ignoreUsersList')">
						<select-twitch-users v-model="formValue.ignoreSettings.users" />
					</n-form-item>
				</n-tab-pane>

				<n-tab-pane name="sounds" :tab="t('overlays.dudes.dudeSoundsDivider')">
					<n-form-item
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<n-switch v-model:value="formValue.dudeSettings.soundsEnabled" />
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.dudeSoundsVolume')">
						<n-slider
							v-model:value="formValue.dudeSettings.soundsVolume"
							:min="0.01"
							:max="1"
							:step="0.01"
							:format-tooltip="(value) => `${(value * 100).toFixed(0)}%`"
							:disabled="!formValue.dudeSettings.soundsEnabled"
						/>
					</n-form-item>
				</n-tab-pane>

				<n-tab-pane name="grow" :tab="t('overlays.dudes.growDivider')">
					<n-form-item :show-feedback="false" :label="t('overlays.dudes.growTime')">
						<n-slider
							v-model:value="formValue.dudeSettings.growTime"
							:min="5000"
							:max="1000 * 60 * 60"
							:step="1000"
							:format-tooltip="(value) => formatDuration(value)"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.growMaxScale')">
						<n-slider
							v-model:value="formValue.dudeSettings.growMaxScale"
							:min="formValue.dudeSettings.scale + 1"
							:max="32"
							:step="1"
						/>
					</n-form-item>
				</n-tab-pane>

				<n-tab-pane name="name-box" :tab="t('overlays.dudes.nameBoxDivider')">
					<n-scrollbar style="max-height: calc(62vh - var(--layout-header-height));" trigger="none">
						<div class="pr-4">
							<n-form>
								<n-form-item
									class="form-item-switch"
									:show-feedback="false"
									:label="t('overlays.dudes.enable')"
								>
									<n-switch v-model:value="formValue.dudeSettings.visibleName" />
								</n-form-item>

								<n-form-item
									:validation-status="nameBoxFillMessage ? 'error' : undefined"
									:feedback="nameBoxFillMessage"
									:label="t('overlays.dudes.nameBoxFill')"
								>
									<n-dynamic-tags
										v-model:value="formValue.nameBoxSettings.fill"
										:disabled="isNameBoxDisabled"
										:max="6"
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
												style="width: 80px;"
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
								</n-form-item>
							</n-form>

							<n-form>
								<n-form-item
									:disabled="isNameBoxDisabled"
									:validation-status="fillGradidentStopMessage ? 'error' : undefined"
									:feedback="fillGradidentStopMessage"
									:label="t('overlays.dudes.nameBoxFillGradientStops')"
								>
									<n-dynamic-tags
										v-model:value="fillGradientStops"
										:render-tag="(tag: string, index: number) => {
											return h(NTag, {
												closable: true,
												onClose: () => {
													formValue.nameBoxSettings.fillGradientStops.splice(index, 1)
												}
											}, { default: () => tag })
										}"
										:max="formValue.nameBoxSettings.fill.length"
										@update:value="(values: string[]) => {
											formValue.nameBoxSettings.fillGradientStops = values.map(Number)
										}"
									>
										<template #input="{ submit, deactivate }">
											<n-input-number
												style="width: 100px;"
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
								</n-form-item>
							</n-form>

							<n-form-item :label="t('overlays.dudes.nameBoxGradientType')">
								<n-select
									v-model:value="formValue.nameBoxSettings.fillGradientType"
									:disabled="isNameBoxDisabled || formValue.nameBoxSettings.fill.length < 2"
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
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxFontFamily')">
								<font-selector
									v-model:font="fontData"
									:disabled="isNameBoxDisabled"
									:font-family="formValue.nameBoxSettings.fontFamily"
									:font-weight="formValue.nameBoxSettings.fontWeight"
									:font-style="formValue.nameBoxSettings.fontStyle"
									:subsets="['latin', 'cyrillic']"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxFontWeight')">
								<n-select
									v-model:value="formValue.nameBoxSettings.fontWeight"
									:disabled="isNameBoxDisabled"
									:options="fontWeightOptions"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxFontStyle')">
								<n-select
									v-model:value="formValue.nameBoxSettings.fontStyle"
									:disabled="isNameBoxDisabled"
									:options="fontStyleOptions"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxFontVariant')">
								<n-select
									v-model:value="formValue.nameBoxSettings.fontVariant"
									:disabled="isNameBoxDisabled"
									:options="fontVariantOptions"
								/>
							</n-form-item>

							<n-form-item :label="t('overlays.dudes.nameBoxFontSize')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.fontSize"
									:disabled="isNameBoxDisabled"
									:min="1"
									:max="128"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxStroke')">
								<n-color-picker
									v-model:value="formValue.nameBoxSettings.stroke"
									:disabled="isNameBoxDisabled"
									:modes="['hex']"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameStrokeThickness')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.strokeThickness"
									:disabled="isNameBoxDisabled"
									:min="0"
									:max="16"
									:step="1"
								/>
							</n-form-item>

							<n-form-item :label="t('overlays.dudes.nameBoxLineJoin')">
								<n-select
									v-model:value="formValue.nameBoxSettings.lineJoin"
									:disabled="isNameBoxDisabled"
									:options="lineJoinOptions"
								/>
							</n-form-item>

							<n-form-item
								class="form-item-switch"
								:show-feedback="false"
								:label="t('overlays.dudes.nameBoxDropShadow')"
							>
								<n-switch
									v-model:value="formValue.nameBoxSettings.dropShadow"
									:disabled="isNameBoxDisabled"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowColor')">
								<n-color-picker
									v-model:value="formValue.nameBoxSettings.dropShadowColor"
									:modes="['hex']"
									:disabled="isDropShadowDisabled"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowAlpha')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.dropShadowAlpha"
									:min="0"
									:max="1"
									:step="0.01"
									:disabled="isDropShadowDisabled"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowBlur')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.dropShadowBlur"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowDistance')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.dropShadowDistance"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
								/>
							</n-form-item>

							<n-form-item :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowAngle')">
								<n-slider
									v-model:value="formValue.nameBoxSettings.dropShadowAngle"
									:min="0"
									:max="Math.PI * 2"
									:step="0.01"
									:format-tooltip="(value) => `${Math.round((value * 180) / Math.PI)}°`"
									:disabled="isDropShadowDisabled"
								/>
							</n-form-item>
						</div>
					</n-scrollbar>
				</n-tab-pane>

				<n-tab-pane name="message-box" :tab="t('overlays.dudes.messageBoxDivider')">
					<n-form-item
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<n-switch v-model:value="formValue.messageBoxSettings.enabled" />
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxShowTime')">
						<n-slider
							v-model:value="formValue.messageBoxSettings.showTime"
							:min="1000"
							:max="60 * 1000"
							:step="1000"
							:format-tooltip="(value) => `${Math.round(value / 1000)}s`"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxFill')">
						<n-color-picker
							v-model:value="formValue.messageBoxSettings.fill"
							:modes="['hex']"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxBackground')">
						<n-color-picker
							v-model:value="formValue.messageBoxSettings.boxColor"
							:modes="['hex']"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxPadding')">
						<n-slider
							v-model:value="formValue.messageBoxSettings.padding"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxBorderRadius')">
						<n-slider
							v-model:value="formValue.messageBoxSettings.borderRadius"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>

					<n-form-item :show-feedback="false" :label="t('overlays.dudes.messageBoxFontSize')">
						<n-slider
							v-model:value="formValue.messageBoxSettings.fontSize"
							:min="12"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</n-form-item>
				</n-tab-pane>

				<n-tab-pane name="emote" :tab="t('overlays.dudes.emoteDivider')">
					<n-form-item
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<n-switch v-model:value="formValue.spitterEmoteSettings.enabled" />
					</n-form-item>
				</n-tab-pane>
			</n-tabs>
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

.form-item-switch {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.card {
	background-color: v-bind('themeVars.cardColor');
}
</style>
