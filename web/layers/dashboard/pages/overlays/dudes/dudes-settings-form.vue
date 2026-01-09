<script setup lang="ts">
import { type Font, FontSelector } from '~/lib/fontsource'
import { DudesSprite } from '@twir/types'
import { addZero, capitalize, colorBrightness, hexToRgb } from '@zero-dependency/utils'
import { intervalToDuration } from 'date-fns'
import {
	MessageSquareIcon,
	MusicIcon,
	PaletteIcon,
	SmileIcon,
	TrendingUpIcon,
	UserIcon,
	UsersIcon,
	XIcon,
} from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import { toast } from 'vue-sonner'

import { useDudesForm } from './use-dudes-form.js'
import { useDudesIframe } from './use-dudes-frame.js'

import { useProfile, useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useDudesOverlayManager } from '#layers/dashboard/api/overlays/dudes'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink.js'
import SelectTwitchUsers from '#layers/dashboard/components/twitchUsers/twitch-users-select.vue'











import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
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

	toast.success(t('sharedTexts.saved'), {
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

// Color picker state for dynamic tags
const showColorPicker = ref(false)
const tempColor = ref('#74f2ca')

function addFillColor() {
	if (formValue.value.nameBoxSettings.fill.length < 6) {
		formValue.value.nameBoxSettings.fill.push(tempColor.value)
		showColorPicker.value = false
		tempColor.value = '#74f2ca'
	}
}

function removeFillColor(index: number) {
	formValue.value.nameBoxSettings.fill.splice(index, 1)
}

// Gradient stop input
const showGradientStopInput = ref(false)
const tempGradientStop = ref(0.1)

function addGradientStop() {
	if (formValue.value.nameBoxSettings.fillGradientStops.length < formValue.value.nameBoxSettings.fill.length) {
		formValue.value.nameBoxSettings.fillGradientStops.push(tempGradientStop.value)
		showGradientStopInput.value = false
		tempGradientStop.value = 0.1
	}
}

function removeGradientStop(index: number) {
	formValue.value.nameBoxSettings.fillGradientStops.splice(index, 1)
}
</script>

<template>
	<UiCard v-if="formValue">
		<UiCardHeader>
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="flex flex-wrap gap-2">
					<UiButton variant="destructive" size="sm" @click="reset">
						{{ t('sharedButtons.setDefaultSettings') }}
					</UiButton>
					<UiButton
						variant="secondary"
						size="sm"
						:disabled="!formValue.id || !canCopyLink"
						@click="copyOverlayLink({ id: formValue.id! })"
					>
						{{ t('overlays.copyOverlayLink') }}
					</UiButton>
					<UiButton size="sm" @click="save">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</div>
			</div>
		</UiCardHeader>

		<UiCardContent>
			<UiAccordion type="multiple" class="w-full" default-value="['dude']" :unmountOnHide="false">
				<!-- Dude Section -->
				<UiAccordionItem value="dude">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<UserIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.dudeDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeDefaultSprite') }}</UiLabel>
					<UiSelect v-model="formValue.dudeSettings.defaultSprite">
						<UiSelectTrigger>
							<UiSelectValue />
						</UiSelectTrigger>
						<UiSelectContent>
							<UiSelectItem v-for="sprite in dudesSprites" :key="sprite.value" :value="sprite.value">
								{{ sprite.label }}
							</UiSelectItem>
						</UiSelectContent>
					</UiSelect>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeMaxOnScreen') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.maxOnScreen]"
							:min="0"
							:max="128"
							:step="1"
							@update:model-value="(val) => formValue.dudeSettings.maxOnScreen = val?.[0] ?? 0"
						/>
						<span class="text-sm text-muted-foreground w-20">
							{{ formValue.dudeSettings.maxOnScreen === 0 ? t('overlays.dudes.dudeMaxOnScreenUnlimited') : formValue.dudeSettings.maxOnScreen }}
						</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeColor') }}</UiLabel>
					<UiColorPicker v-model="formValue.dudeSettings.color" />
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeGravity') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.gravity]"
							:min="100"
							:max="5000"
							@update:model-value="(val) => formValue.dudeSettings.gravity = val?.[0] ?? 100"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.dudeSettings.gravity }}</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeMaxLifeTime') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.maxLifeTime]"
							:min="1000"
							:max="120 * 60 * 1000"
							:step="1000"
							@update:model-value="(val) => formValue.dudeSettings.maxLifeTime = val?.[0] ?? 1000"
						/>
						<span class="text-sm text-muted-foreground w-32">{{ formatDuration(formValue.dudeSettings.maxLifeTime) }}</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeScale') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.scale]"
							:min="1"
							:max="10"
							:step="1"
							@update:model-value="(val) => formValue.dudeSettings.scale = val?.[0] ?? 1"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.dudeSettings.scale }}</span>
					</div>
				</div>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Ignoring Section -->
				<UiAccordionItem value="ignoring">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<UsersIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.ignoreDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<div class="flex items-center justify-between">
					<UiLabel>{{ t('overlays.dudes.ignoreCommands') }}</UiLabel>
					<UiSwitch
						:model-value="formValue.ignoreSettings.ignoreCommands"
						@update:model-value="formValue.ignoreSettings.ignoreCommands = $event"
					/>
				</div>

				<div class="flex items-center justify-between">
					<UiLabel>{{ t('overlays.dudes.ignoreUsers') }}</UiLabel>
					<UiSwitch
						:model-value="formValue.ignoreSettings.ignoreUsers"
						@update:model-value="formValue.ignoreSettings.ignoreUsers = $event"
					/>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.ignoreUsersList') }}</UiLabel>
					<SelectTwitchUsers v-model="formValue.ignoreSettings.users" />
				</div>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Sounds Section -->
				<UiAccordionItem value="sounds">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<MusicIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.dudeSoundsDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<div class="flex items-center justify-between">
					<UiLabel>{{ t('overlays.dudes.enable') }}</UiLabel>
					<UiSwitch
						:model-value="formValue.dudeSettings.soundsEnabled"
						@update:model-value="formValue.dudeSettings.soundsEnabled = $event"
					/>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.dudeSoundsVolume') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.soundsVolume]"
							:min="0.01"
							:max="1"
							:step="0.01"
							:disabled="!formValue.dudeSettings.soundsEnabled"
							@update:model-value="(val) => formValue.dudeSettings.soundsVolume = val?.[0] ?? 0.01"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ (formValue.dudeSettings.soundsVolume * 100).toFixed(0) }}%</span>
					</div>
				</div>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Grow Section -->
				<UiAccordionItem value="grow">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<TrendingUpIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.growDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.growTime') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.growTime]"
							:min="5000"
							:max="1000 * 60 * 60"
							:step="1000"
							@update:model-value="(val) => formValue.dudeSettings.growTime = val?.[0] ?? 5000"
						/>
						<span class="text-sm text-muted-foreground w-32">{{ formatDuration(formValue.dudeSettings.growTime) }}</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.growMaxScale') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.dudeSettings.growMaxScale]"
							:min="formValue.dudeSettings.scale + 1"
							:max="32"
							:step="1"
							@update:model-value="(val) => formValue.dudeSettings.growMaxScale = val?.[0] ?? 1"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.dudeSettings.growMaxScale }}</span>
					</div>
				</div>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Name Box Section -->
				<UiAccordionItem value="name-box">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<PaletteIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.nameBoxDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<UiScrollArea class="h-[50vh] pr-4">
					<div class="space-y-4">
						<div class="flex items-center justify-between">
							<UiLabel>{{ t('overlays.dudes.enable') }}</UiLabel>
							<UiSwitch
								:model-value="formValue.dudeSettings.visibleName"
								@update:model-value="formValue.dudeSettings.visibleName = $event"
							/>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFill') }}</UiLabel>
							<div class="flex flex-wrap gap-2">
								<UiBadge
									v-for="(color, index) in formValue.nameBoxSettings.fill"
									:key="index"
									class="pr-1"
									:style="{ backgroundColor: color, color: (hexToRgb(color) && colorBrightness(hexToRgb(color)!) > 128) ? '#000' : '#fff' }"
								>
									{{ color }}
									<UiButton
										variant="ghost"
										size="icon"
										class="h-4 w-4 ml-1 hover:bg-transparent"
										:disabled="isNameBoxDisabled"
										@click="removeFillColor(index)"
									>
										<XIcon class="h-3 w-3" />
									</UiButton>
								</UiBadge>
								<UiButton
									v-if="!showColorPicker && formValue.nameBoxSettings.fill.length < 6"
									variant="outline"
									size="sm"
									:disabled="isNameBoxDisabled"
									@click="showColorPicker = true"
								>
									+ Add Color
								</UiButton>
								<div v-if="showColorPicker" class="flex items-center gap-2">
									<UiColorPicker v-model="tempColor" class="w-20" />
									<UiButton size="sm" @click="addFillColor">Add</UiButton>
									<UiButton size="sm" variant="ghost" @click="showColorPicker = false">Cancel</UiButton>
								</div>
							</div>
							<span v-if="nameBoxFillMessage" class="text-sm text-destructive">{{ nameBoxFillMessage }}</span>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFillGradientStops') }}</UiLabel>
							<div class="flex flex-wrap gap-2">
								<UiBadge
									v-for="(stop, index) in formValue.nameBoxSettings.fillGradientStops"
									:key="index"
									variant="secondary"
									class="pr-1"
								>
									{{ stop }}
									<UiButton
										variant="ghost"
										size="icon"
										class="h-4 w-4 ml-1 hover:bg-transparent"
										:disabled="isNameBoxDisabled"
										@click="removeGradientStop(index)"
									>
										<XIcon class="h-3 w-3" />
									</UiButton>
								</UiBadge>
								<UiButton
									v-if="!showGradientStopInput && formValue.nameBoxSettings.fillGradientStops.length < formValue.nameBoxSettings.fill.length"
									variant="outline"
									size="sm"
									:disabled="isNameBoxDisabled"
									@click="showGradientStopInput = true"
								>
									+ Add Stop
								</UiButton>
								<div v-if="showGradientStopInput" class="flex items-center gap-2">
									<UiInput
										v-model.number="tempGradientStop"
										type="number"
										:min="0"
										:max="1"
										:step="0.01"
										class="w-24"
									/>
									<UiButton size="sm" @click="addGradientStop">Add</UiButton>
									<UiButton size="sm" variant="ghost" @click="showGradientStopInput = false">Cancel</UiButton>
								</div>
							</div>
							<span v-if="fillGradidentStopMessage" class="text-sm text-destructive">{{ fillGradidentStopMessage }}</span>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxGradientType') }}</UiLabel>
							<UiSelect
								v-model="formValue.nameBoxSettings.fillGradientType"
								:disabled="isNameBoxDisabled || formValue.nameBoxSettings.fill.length < 2"
							>
								<UiSelectTrigger>
									<UiSelectValue />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem :value="0">Vertical</UiSelectItem>
									<UiSelectItem :value="1">Horizontal</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFontFamily') }}</UiLabel>
							<FontSelector
								v-model:font="fontData"
								:disabled="isNameBoxDisabled"
								:font-family="formValue.nameBoxSettings.fontFamily"
								:font-weight="formValue.nameBoxSettings.fontWeight"
								:font-style="formValue.nameBoxSettings.fontStyle"
								:subsets="['latin', 'cyrillic']"
							/>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFontWeight') }}</UiLabel>
							<UiSelect v-model="formValue.nameBoxSettings.fontWeight" :disabled="isNameBoxDisabled">
								<UiSelectTrigger>
									<UiSelectValue />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem v-for="option in fontWeightOptions" :key="option.value" :value="option.value">
										{{ option.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFontStyle') }}</UiLabel>
							<UiSelect v-model="formValue.nameBoxSettings.fontStyle" :disabled="isNameBoxDisabled">
								<UiSelectTrigger>
									<UiSelectValue />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem v-for="option in fontStyleOptions" :key="option.value" :value="option.value">
										{{ option.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFontVariant') }}</UiLabel>
							<UiSelect v-model="formValue.nameBoxSettings.fontVariant" :disabled="isNameBoxDisabled">
								<UiSelectTrigger>
									<UiSelectValue />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem v-for="option in fontVariantOptions" :key="option.value" :value="option.value">
										{{ option.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxFontSize') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.fontSize]"
									:min="1"
									:max="128"
									:disabled="isNameBoxDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.fontSize = val?.[0] ?? 1"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ formValue.nameBoxSettings.fontSize }}</span>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxStroke') }}</UiLabel>
							<UiColorPicker v-model="formValue.nameBoxSettings.stroke" :disabled="isNameBoxDisabled" />
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameStrokeThickness') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.strokeThickness]"
									:min="0"
									:max="16"
									:step="1"
									:disabled="isNameBoxDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.strokeThickness = val?.[0] ?? 0"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ formValue.nameBoxSettings.strokeThickness }}</span>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxLineJoin') }}</UiLabel>
							<UiSelect v-model="formValue.nameBoxSettings.lineJoin" :disabled="isNameBoxDisabled">
								<UiSelectTrigger>
									<UiSelectValue />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem v-for="option in lineJoinOptions" :key="option.value" :value="option.value">
										{{ option.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</div>

						<div class="flex items-center justify-between">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadow') }}</UiLabel>
							<UiSwitch
								:model-value="formValue.nameBoxSettings.dropShadow"
								:disabled="isNameBoxDisabled"
								@update:model-value="formValue.nameBoxSettings.dropShadow = $event"
							/>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadowColor') }}</UiLabel>
							<UiColorPicker v-model="formValue.nameBoxSettings.dropShadowColor" :disabled="isDropShadowDisabled" />
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadowAlpha') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.dropShadowAlpha]"
									:min="0"
									:max="1"
									:step="0.01"
									:disabled="isDropShadowDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.dropShadowAlpha = val?.[0] ?? 0"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ formValue.nameBoxSettings.dropShadowAlpha.toFixed(2) }}</span>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadowBlur') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.dropShadowBlur]"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.dropShadowBlur = val?.[0] ?? 0"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ formValue.nameBoxSettings.dropShadowBlur.toFixed(1) }}</span>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadowDistance') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.dropShadowDistance]"
									:min="0"
									:max="32"
									:step="0.1"
									:disabled="isDropShadowDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.dropShadowDistance = val?.[0] ?? 0"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ formValue.nameBoxSettings.dropShadowDistance.toFixed(1) }}</span>
							</div>
						</div>

						<div class="flex flex-col gap-2">
							<UiLabel>{{ t('overlays.dudes.nameBoxDropShadowAngle') }}</UiLabel>
							<div class="flex items-center gap-4">
								<UiSlider
									:model-value="[formValue.nameBoxSettings.dropShadowAngle]"
									:min="0"
									:max="Math.PI * 2"
									:step="0.01"
									:disabled="isDropShadowDisabled"
									@update:model-value="(val) => formValue.nameBoxSettings.dropShadowAngle = val?.[0] ?? 0"
								/>
								<span class="text-sm text-muted-foreground w-20">{{ Math.round((formValue.nameBoxSettings.dropShadowAngle * 180) / Math.PI) }}°</span>
							</div>
						</div>
					</div>
				</UiScrollArea>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Message Box Section -->
				<UiAccordionItem value="message-box">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<MessageSquareIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.messageBoxDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
				<div class="flex items-center justify-between">
					<UiLabel>{{ t('overlays.dudes.enable') }}</UiLabel>
					<UiSwitch
						:model-value="formValue.messageBoxSettings.enabled"
						@update:model-value="formValue.messageBoxSettings.enabled = $event"
					/>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxShowTime') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.messageBoxSettings.showTime]"
							:min="1000"
							:max="60 * 1000"
							:step="1000"
							:disabled="isMessageBoxDisabled"
							@update:model-value="(val) => formValue.messageBoxSettings.showTime = val?.[0] ?? 1000"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ Math.round(formValue.messageBoxSettings.showTime / 1000) }}s</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxFill') }}</UiLabel>
					<UiColorPicker v-model="formValue.messageBoxSettings.fill" :disabled="isMessageBoxDisabled" />
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxBackground') }}</UiLabel>
					<UiColorPicker v-model="formValue.messageBoxSettings.boxColor" :disabled="isMessageBoxDisabled" />
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxPadding') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.messageBoxSettings.padding]"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
							@update:model-value="(val) => formValue.messageBoxSettings.padding = val?.[0] ?? 0"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.messageBoxSettings.padding }}</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxBorderRadius') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.messageBoxSettings.borderRadius]"
							:min="0"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
							@update:model-value="(val) => formValue.messageBoxSettings.borderRadius = val?.[0] ?? 0"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.messageBoxSettings.borderRadius }}</span>
					</div>
				</div>

				<div class="flex flex-col gap-2">
					<UiLabel>{{ t('overlays.dudes.messageBoxFontSize') }}</UiLabel>
					<div class="flex items-center gap-4">
						<UiSlider
							:model-value="[formValue.messageBoxSettings.fontSize]"
							:min="12"
							:max="64"
							:step="1"
							:disabled="isMessageBoxDisabled"
							@update:model-value="(val) => formValue.messageBoxSettings.fontSize = val?.[0] ?? 12"
						/>
						<span class="text-sm text-muted-foreground w-20">{{ formValue.messageBoxSettings.fontSize }}</span>
					</div>
				</div>
					</UiAccordionContent>
				</UiAccordionItem>

				<!-- Emote Section -->
				<UiAccordionItem value="emote">
					<UiAccordionTrigger>
						<div class="flex items-center gap-2">
							<SmileIcon class="h-4 w-4" />
							<span>{{ t('overlays.dudes.emoteDivider') }}</span>
						</div>
					</UiAccordionTrigger>
					<UiAccordionContent class="space-y-4 pt-4">
						<div class="flex items-center justify-between">
							<UiLabel>{{ t('overlays.dudes.enable') }}</UiLabel>
							<UiSwitch
								:model-value="formValue.spitterEmoteSettings.enabled"
								@update:model-value="formValue.spitterEmoteSettings.enabled = $event"
							/>
						</div>
					</UiAccordionContent>
				</UiAccordionItem>
			</UiAccordion>
		</UiCardContent>
	</UiCard>
</template>

<style scoped>
@import '../styles.css';
</style>
