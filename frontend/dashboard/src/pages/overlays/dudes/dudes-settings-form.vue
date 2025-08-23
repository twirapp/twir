<script setup lang="ts">
import { FontSelector } from '@twir/fontsource'
import { DudesSprite } from '@twir/types'
import { addZero, capitalize, colorBrightness, hexToRgb } from '@zero-dependency/utils'
import { intervalToDuration } from 'date-fns'
import {
	NButton,
	NColorPicker,
	NDynamicTags,
	NForm,
	NFormItem,
	NInputNumber,
	NScrollbar,
	NSelect,
	NSlider,
	NSwitch,
	NTabPane,
	NTabs,
	NTag,
	useThemeVars,
} from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useDudesForm } from './use-dudes-form.js'
import { useDudesIframe } from './use-dudes-frame.js'

import type { Font } from '@twir/fontsource'

import { useDudesOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api/index.js'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js'
import SelectTwitchUsers from '@/components/twitchUsers/multiple.vue'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const themeVars = useThemeVars()
const discrete = useNaiveDiscrete()
const { copyOverlayLink } = useCopyOverlayLink('dudes')
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { data: profile } = useProfile()
const { data: formValue, reset } = useDudesForm()
const { sendIframeMessage } = useDudesIframe()

watch(
	formValue,
	(form) => {
		if (!form) return
		if (!form.nameBoxSettings.fill.length) return
		sendIframeMessage('update-settings', form)
	},
	{ deep: true }
)

watch(
	() => formValue.value.dudeSettings.defaultSprite,
	(dudeSprite) => {
		sendIframeMessage('update-sprite', dudeSprite)
	}
)

watch(
	() => formValue.value.dudeSettings.color,
	(dudeColor) => {
		sendIframeMessage('update-color', dudeColor)
	}
)

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays
})

const manager = useDudesOverlayManager()
const updater = manager.useUpdate()

async function save() {
	if (!formValue.value.id) return

	// Extract the settings without the id for the input
	const { id, ...settings } = formValue.value

	await updater.executeMutation({
		id: formValue.value.id,
		input: settings,
	})

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	})
}

function formatDuration(duration: number) {
	const { hours = 0, minutes = 0, seconds = 0 } = intervalToDuration({ start: 0, end: duration })
	if (hours === 0) {
		return `${addZero(minutes)}:${addZero(seconds)}`
	}

	return `${addZero(hours)}:${addZero(minutes)}:${addZero(seconds)}`
}

const fillGradientStops = computed(() => {
	if (!formValue.value) return []
	return formValue.value.nameBoxSettings.fillGradientStops.map((stop) => `${stop}`)
})

const fillGradidentStopMessage = computed(() => {
	if (!formValue.value.nameBoxSettings.fillGradientStops.length) {
		return t('overlays.dudes.nameBoxFillGradientStopsError')
	}

	return ''
})

const nameBoxFillMessage = computed(() => {
	if (!formValue.value.nameBoxSettings.fill.length) {
		return t('overlays.dudes.nameBoxFillError')
	}

	return ''
})

const fontData = ref<Font | null>(null)
watch(
	() => fontData.value,
	(font) => {
		if (!font) return
		formValue.value.nameBoxSettings.fontFamily = font.id
		formValue.value.messageBoxSettings.fontFamily = font.id
	}
)

const fontWeightOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.weights.map((weight) => ({
		label: `${weight}`,
		value: weight,
	}))
})

const fontStyleOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.styles.map((style) => ({
		label: capitalize(style),
		value: style,
	}))
})

const fontVariantOptions = ['normal', 'small-caps'].map((variant) => ({
	label: capitalize(variant),
	value: variant,
}))

const lineJoinOptions = ['round', 'bevel', 'miter'].map((lineJoin) => ({
	label: capitalize(lineJoin),
	value: lineJoin,
}))

const isMessageBoxDisabled = computed(() => {
	return !formValue.value.messageBoxSettings.enabled
})

const isNameBoxDisabled = computed(() => {
	return !formValue.value.dudeSettings.visibleName
})

const isDropShadowDisabled = computed(() => {
	return isNameBoxDisabled.value || !formValue.value.nameBoxSettings.dropShadow
})

const dudesSprites = Object.keys(DudesSprite).map((key) => ({
	label: capitalize(key),
	value: key,
}))
</script>

<template>
	<div v-if="formValue" class="card">
		<div class="card-header">
			<NButton secondary type="error" @click="reset">
				{{ t('sharedButtons.setDefaultSettings') }}
			</NButton>
			<NButton
				secondary
				type="info"
				:disabled="!formValue.id || !canCopyLink"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</NButton>
			<NButton secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</NButton>

			<NTabs class="pt-2" type="line" placement="left" default-value="dude" animated>
				<NTabPane name="dude" :tab="t('overlays.dudes.dudeDivider')">
					<NFormItem :label="t('overlays.dudes.dudeDefaultSprite')">
						<NSelect v-model:value="formValue.dudeSettings.defaultSprite" :options="dudesSprites" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.dudeMaxOnScreen')">
						<NSlider
							v-model:value="formValue.dudeSettings.maxOnScreen"
							:min="0"
							:max="128"
							:step="1"
							:format-tooltip="
								(value) => {
									if (value === 0) {
										return t('overlays.dudes.dudeMaxOnScreenUnlimited')
									}

									return value
								}
							"
						/>
					</NFormItem>

					<NFormItem :label="t('overlays.dudes.dudeColor')">
						<NColorPicker v-model:value="formValue.dudeSettings.color" :modes="['hex']" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.dudeGravity')">
						<NSlider v-model:value="formValue.dudeSettings.gravity" :min="100" :max="5000" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.dudeMaxLifeTime')">
						<NSlider
							v-model:value="formValue.dudeSettings.maxLifeTime"
							:min="1000"
							:max="120 * 60 * 1000"
							:step="1000"
							:format-tooltip="(value) => formatDuration(value)"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.dudeScale')">
						<NSlider v-model:value="formValue.dudeSettings.scale" :min="1" :max="10" :step="1" />
					</NFormItem>
				</NTabPane>

				<NTabPane name="ignoring" :tab="t('overlays.dudes.ignoreDivider')">
					<NFormItem
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.ignoreCommands')"
					>
						<NSwitch v-model:value="formValue.ignoreSettings.ignoreCommands" />
					</NFormItem>

					<NFormItem
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.ignoreUsers')"
					>
						<NSwitch v-model:value="formValue.ignoreSettings.ignoreUsers" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.ignoreUsersList')">
						<SelectTwitchUsers v-model="formValue.ignoreSettings.users" />
					</NFormItem>
				</NTabPane>

				<NTabPane name="sounds" :tab="t('overlays.dudes.dudeSoundsDivider')">
					<NFormItem
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<NSwitch v-model:value="formValue.dudeSettings.soundsEnabled" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.dudeSoundsVolume')">
						<NSlider
							v-model:value="formValue.dudeSettings.soundsVolume"
							:min="0.01"
							:max="1"
							:step="0.01"
							:format-tooltip="(value) => `${(value * 100).toFixed(0)}%`"
							:disabled="!formValue.dudeSettings.soundsEnabled"
						/>
					</NFormItem>
				</NTabPane>

				<NTabPane name="grow" :tab="t('overlays.dudes.growDivider')">
					<NFormItem :show-feedback="false" :label="t('overlays.dudes.growTime')">
						<NSlider
							v-model:value="formValue.dudeSettings.growTime"
							:min="5000"
							:max="1000 * 60 * 60"
							:step="1000"
							:format-tooltip="(value) => formatDuration(value)"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.growMaxScale')">
						<NSlider
							v-model:value="formValue.dudeSettings.growMaxScale"
							:min="formValue.dudeSettings.scale + 1"
							:max="32"
							:step="1"
						/>
					</NFormItem>
				</NTabPane>

				<NTabPane name="name-box" :tab="t('overlays.dudes.nameBoxDivider')">
					<NScrollbar style="max-height: calc(62vh - var(--layout-header-height))" trigger="none">
						<div class="pr-4">
							<NForm>
								<NFormItem
									class="form-item-switch"
									:show-feedback="false"
									:label="t('overlays.dudes.enable')"
								>
									<NSwitch v-model:value="formValue.dudeSettings.visibleName" />
								</NFormItem>

								<NFormItem
									:validation-status="nameBoxFillMessage ? 'error' : undefined"
									:feedback="nameBoxFillMessage"
									:label="t('overlays.dudes.nameBoxFill')"
								>
									<NDynamicTags
										v-model:value="formValue.nameBoxSettings.fill"
										:disabled="isNameBoxDisabled"
										:max="6"
										:render-tag="
											(tag: string, index: number) => {
												const rgb = hexToRgb(tag)
												const textColor = rgb && colorBrightness(rgb) > 128 ? '#000' : '#fff'

												return h(
													NTag,
													{
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
														},
													},
													{ default: () => tag }
												)
											}
										"
									>
										<template #input="{ submit, deactivate }">
											<NColorPicker
												style="width: 80px"
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
									</NDynamicTags>
								</NFormItem>
							</NForm>

							<NForm>
								<NFormItem
									:disabled="isNameBoxDisabled"
									:validation-status="fillGradidentStopMessage ? 'error' : undefined"
									:feedback="fillGradidentStopMessage"
									:label="t('overlays.dudes.nameBoxFillGradientStops')"
								>
									<NDynamicTags
										v-model:value="fillGradientStops"
										:render-tag="
											(tag: string, index: number) => {
												return h(
													NTag,
													{
														closable: true,
														onClose: () => {
															formValue.nameBoxSettings.fillGradientStops.splice(index, 1)
														},
													},
													{ default: () => tag }
												)
											}
										"
										:max="formValue.nameBoxSettings.fill.length"
										@update:value="
											(values: string[]) => {
												formValue.nameBoxSettings.fillGradientStops = values.map(Number)
											}
										"
									>
										<template #input="{ submit, deactivate }">
											<NInputNumber
												style="width: 100px"
												autofocus
												placeholder=""
												:max="1"
												:min="0"
												:step="0.01"
												:default-value="0.1"
												size="small"
												:update-value-on-input="false"
												:parse="
													(v) => {
														const parsedNum = Number(v)
														return Number.isNaN(parsedNum) ? 0 : parsedNum
													}
												"
												@keyup.enter="submit($event.target.value)"
												@confirm="submit($event)"
												@blur="deactivate"
											/>
										</template>
									</NDynamicTags>
								</NFormItem>
							</NForm>

							<NFormItem :label="t('overlays.dudes.nameBoxGradientType')">
								<NSelect
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
										},
									]"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxFontFamily')">
								<FontSelector
									v-model:font="fontData"
									:disabled="isNameBoxDisabled"
									:font-family="formValue.nameBoxSettings.fontFamily"
									:font-weight="formValue.nameBoxSettings.fontWeight"
									:font-style="formValue.nameBoxSettings.fontStyle"
									:subsets="['latin', 'cyrillic']"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxFontWeight')">
								<NSelect
									v-model:value="formValue.nameBoxSettings.fontWeight"
									:disabled="isNameBoxDisabled"
									:options="fontWeightOptions"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxFontStyle')">
								<NSelect
									v-model:value="formValue.nameBoxSettings.fontStyle"
									:disabled="isNameBoxDisabled"
									:options="fontStyleOptions"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxFontVariant')">
								<NSelect
									v-model:value="formValue.nameBoxSettings.fontVariant"
									:disabled="isNameBoxDisabled"
									:options="fontVariantOptions"
								/>
							</NFormItem>

							<NFormItem :label="t('overlays.dudes.nameBoxFontSize')">
								<NSlider
									v-model:value="formValue.nameBoxSettings.fontSize"
									:disabled="isNameBoxDisabled"
									:min="1"
									:max="128"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxStroke')">
								<NColorPicker
									v-model:value="formValue.nameBoxSettings.stroke"
									:disabled="isNameBoxDisabled"
									:modes="['hex']"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameStrokeThickness')">
								<NSlider
									v-model:value="formValue.nameBoxSettings.strokeThickness"
									:disabled="isNameBoxDisabled"
									:min="0"
									:max="16"
									:step="1"
								/>
							</NFormItem>

							<NFormItem :label="t('overlays.dudes.nameBoxLineJoin')">
								<NSelect
									v-model:value="formValue.nameBoxSettings.lineJoin"
									:disabled="isNameBoxDisabled"
									:options="lineJoinOptions"
								/>
							</NFormItem>

							<NFormItem
								class="form-item-switch"
								:show-feedback="false"
								:label="t('overlays.dudes.nameBoxDropShadow')"
							>
								<NSwitch
									v-model:value="formValue.nameBoxSettings.dropShadow"
									:disabled="isNameBoxDisabled"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowColor')">
								<NColorPicker
									v-model:value="formValue.nameBoxSettings.dropShadowColor"
									:modes="['hex']"
									:disabled="isDropShadowDisabled"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowAlpha')">
								<NSlider
									v-model:value="formValue.nameBoxSettings.dropShadowAlpha"
									:min="0"
									:max="1"
									:step="0.01"
									:disabled="isDropShadowDisabled"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowBlur')">
								<NSlider
									v-model:value="formValue.nameBoxSettings.dropShadowBlur"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
								/>
							</NFormItem>

							<NFormItem
								:show-feedback="false"
								:label="t('overlays.dudes.nameBoxDropShadowDistance')"
							>
								<NSlider
									v-model:value="formValue.nameBoxSettings.dropShadowDistance"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
								/>
							</NFormItem>

							<NFormItem :show-feedback="false" :label="t('overlays.dudes.nameBoxDropShadowAngle')">
								<NSlider
									v-model:value="formValue.nameBoxSettings.dropShadowAngle"
									:min="0"
									:max="Math.PI * 2"
									:step="0.01"
									:format-tooltip="(value) => `${Math.round((value * 180) / Math.PI)}°`"
									:disabled="isDropShadowDisabled"
								/>
							</NFormItem>
						</div>
					</NScrollbar>
				</NTabPane>

				<NTabPane name="message-box" :tab="t('overlays.dudes.messageBoxDivider')">
					<NFormItem
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<NSwitch v-model:value="formValue.messageBoxSettings.enabled" />
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxShowTime')">
						<NSlider
							v-model:value="formValue.messageBoxSettings.showTime"
							:min="1000"
							:max="60 * 1000"
							:step="1000"
							:format-tooltip="(value) => `${Math.round(value / 1000)}s`"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxFill')">
						<NColorPicker
							v-model:value="formValue.messageBoxSettings.fill"
							:modes="['hex']"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxBackground')">
						<NColorPicker
							v-model:value="formValue.messageBoxSettings.boxColor"
							:modes="['hex']"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxPadding')">
						<NSlider
							v-model:value="formValue.messageBoxSettings.padding"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxBorderRadius')">
						<NSlider
							v-model:value="formValue.messageBoxSettings.borderRadius"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>

					<NFormItem :show-feedback="false" :label="t('overlays.dudes.messageBoxFontSize')">
						<NSlider
							v-model:value="formValue.messageBoxSettings.fontSize"
							:min="12"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
						/>
					</NFormItem>
				</NTabPane>

				<NTabPane name="emote" :tab="t('overlays.dudes.emoteDivider')">
					<NFormItem
						class="form-item-switch"
						:show-feedback="false"
						:label="t('overlays.dudes.enable')"
					>
						<NSwitch v-model:value="formValue.spitterEmoteSettings.enabled" />
					</NFormItem>
				</NTabPane>
			</NTabs>
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
