<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { type Font, FontSelector } from '~/lib/fontsource'
import { InfoIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, ref, toRaw, watch } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useBeRightBackOverlayApi } from '#layers/dashboard/api/overlays-be-right-back'









import { toast } from 'vue-sonner'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink.ts'
import CommandButton from '~/features/commands/ui/command-button.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'

import { BeRightBackUpdateInputSchema } from '~/gql/validation-schemas.ts'

const { t } = useI18n()

const { user: profile } = storeToRefs(useDashboardAuth())

const defaultSettings = {
	backgroundColor: 'rgba(9, 8, 8, 0.50)',
	fontColor: '#fff',
	fontFamily: 'inter',
	fontSize: 100,
	text: 'AFK FOR',
	late: {
		text: 'LATE FOR',
		displayBrbTime: true,
		enabled: true,
	},
}

const brbForm = useForm({
	validationSchema: toTypedSchema(BeRightBackUpdateInputSchema),
	initialValues: defaultSettings,
	keepValuesOnUnmount: true,
})

const api = useBeRightBackOverlayApi()
const {
	data: settings,
	fetching: isSettingsLoading,
	error: settingsError,
} = api.useQueryBeRightBack()
const updateMutation = api.useMutationUpdateBeRightBack()

const isSettingsError = computed(() => !!settingsError.value)

watch(
	settings,
	(v) => {
		if (!v?.overlaysBeRightBack) return
		const overlay = v.overlaysBeRightBack
		brbForm.setValues({
			text: overlay.text,
			backgroundColor: overlay.backgroundColor,
			fontColor: overlay.fontColor,
			fontFamily: overlay.fontFamily,
			fontSize: overlay.fontSize,
			late: {
				enabled: overlay.late.enabled,
				text: overlay.late.text,
				displayBrbTime: overlay.late.displayBrbTime,
			},
		})
	},
	{ immediate: true }
)

const brbIframeRef = ref<HTMLIFrameElement | null>(null)
const brbIframeUrl = computed(() => {
	if (!profile.value) return null

	return `${window.location.origin}/overlays/${profile.value.apiKey}/brb`
})

function sendIframeMessage(key: string, data?: any) {
	if (!brbIframeRef.value) return
	const win = brbIframeRef.value

	win.contentWindow?.postMessage(
		JSON.stringify({
			key,
			data: toRaw(data),
		})
	)
}

function sendSettings() {
	const values = brbForm.values
	sendIframeMessage('settings', {
		...toRaw(values),
		channelName: profile.value?.login,
		channelId: profile.value?.id,
	})
}

watch(brbIframeRef, (v) => {
	if (!v) return

	v.contentWindow?.addEventListener('message', (e) => {
		const parsed = JSON.parse(e.data)
		if (parsed.key !== 'getSettings') return

		sendSettings()
	})
})

watch(
	() => brbForm.values,
	() => {
		if (!brbIframeRef.value) return

		sendSettings()
	},
	{ deep: true }
)

const { copyOverlayLink } = useCopyOverlayLink('brb')

const save = brbForm.handleSubmit(async (values) => {
	try {
		await updateMutation.executeMutation({
			input: {
				text: values.text,
				backgroundColor: values.backgroundColor,
				fontColor: values.fontColor,
				fontFamily: values.fontFamily,
				fontSize: values.fontSize,
				late: {
					enabled: values.late.enabled,
					text: values.late.text,
					displayBrbTime: values.late.displayBrbTime,
				},
			},
		})

		toast.success(t('sharedTexts.saved'))
	} catch (e) {
		toast.error('Error occurred while saving BRB overlay')
		console.error(e)
	}
})

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays
})

function setDefaultSettings() {
	brbForm.setValues(structuredClone(defaultSettings))
}

const fontData = ref<Font | null>(null)
watch(
	() => fontData.value,
	(font) => {
		if (!font) return
		brbForm.setFieldValue('fontFamily', font.id)
	}
)
</script>

<template>
	<div class="flex flex-col lg:flex-row gap-4 p-4">
		<UiCard class="flex-1">
			<UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<h2 class="text-2xl font-bold">
					{{ t('overlays.brb.title', 'Be Right Back Overlay') }}
				</h2>
				<div class="flex gap-2">
					<UiButton variant="outline" @click="setDefaultSettings">
						{{ t('sharedButtons.setDefaultSettings') }}
					</UiButton>
					<UiButton
						variant="outline"
						:disabled="isSettingsError || isSettingsLoading || !canCopyLink"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</UiButton>
					<UiButton @click="save">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</div>
			</UiCardHeader>

			<UiCardContent class="space-y-6">
				<form @submit="save" class="space-y-6">
					<div class="space-y-4">
						<UiSeparator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.brb.settings.main.label') }}
						</h3>

						<div class="space-y-4">
							<div class="flex flex-col gap-2">
								<CommandButton
									name="brb"
									:title="t('overlays.brb.settings.main.startCommand.description')"
								/>
								<UiAlert>
									<InfoIcon class="h-4 w-4" />
									<UiAlertDescription>
										<span v-html="t('overlays.brb.settings.main.startCommand.example')" />
									</UiAlertDescription>
								</UiAlert>
								<CommandButton
									name="brbstop"
									:title="t('overlays.brb.settings.main.stopCommand.description')"
								/>
							</div>

							<UiFormField v-slot="{ componentField }" name="text">
								<UiFormItem>
									<UiFormLabel>{{ t('overlays.brb.settings.main.text') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" :maxlength="500" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="backgroundColor">
								<UiFormItem>
									<UiFormLabel>{{ t('overlays.brb.settings.main.background') }}</UiFormLabel>
									<UiFormControl>
										<UiInputWithIcon
											id="backgroundColor"
											v-model="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										>
											<UiColorPicker
												output-format="rgb"
												v-model="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</UiInputWithIcon>
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="fontColor">
								<UiFormItem>
									<UiFormLabel>{{ t('overlays.brb.settings.main.font.color') }}</UiFormLabel>
									<UiFormControl>
										<UiInputWithIcon
											id="fontColor"
											v-model="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										>
											<UiColorPicker
												v-model="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</UiInputWithIcon>
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<div class="space-y-2">
								<UiLabel>{{ t('overlays.brb.settings.main.font.family') }}</UiLabel>
								<FontSelector
									v-model:font="fontData"
									:font-family="brbForm.values.fontFamily || 'inter'"
									font-style="normal"
									:font-weight="400"
								/>
							</div>

							<UiFormField v-slot="{ componentField }" name="fontSize">
								<UiFormItem>
									<UiFormLabel>{{ t('overlays.brb.settings.main.font.size') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" type="number" :min="1" :max="500" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>

						<UiSeparator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.brb.settings.late.label') }}
						</h3>

						<div class="space-y-4">
							<UiFormField v-slot="{ componentField }" name="late.text">
								<UiFormItem>
									<UiFormLabel>{{ t('overlays.brb.settings.late.text') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" :maxlength="500" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ field }" name="late.enabled">
								<UiFormItem>
									<div class="flex items-center gap-2">
										<UiFormControl>
											<UiSwitch
												:model-value="field.value"
												@update:model-value="field['onUpdate:modelValue']"
											/>
										</UiFormControl>
										<UiFormLabel class="mt-0!">{{ t('sharedTexts.enabled') }}</UiFormLabel>
									</div>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ field }" name="late.displayBrbTime">
								<UiFormItem>
									<div class="flex items-center gap-2">
										<UiFormControl>
											<UiSwitch
												:model-value="field.value"
												@update:model-value="field['onUpdate:modelValue']"
											/>
										</UiFormControl>
										<UiFormLabel class="mt-0!">{{
											t('overlays.brb.settings.late.displayBrb')
										}}</UiFormLabel>
									</div>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>
					</div>
				</form>
			</UiCardContent>
		</UiCard>

		<div class="flex-1 relative">
			<iframe
				v-if="brbIframeUrl"
				ref="brbIframeRef"
				:src="brbIframeUrl"
				class="w-full h-[600px] border rounded-lg"
			/>
			<div class="absolute top-4 right-4 flex gap-2">
				<UiButton variant="outline" size="sm" @click="sendIframeMessage('stop')">
					{{ t('overlays.brb.preview.stop') }}
				</UiButton>
				<UiButton
					size="sm"
					@click="
						() => {
							sendSettings()
							sendIframeMessage('start', { minutes: 0.1 })
						}
					"
				>
					{{ t('overlays.brb.preview.start') }}
				</UiButton>
			</div>
		</div>
	</div>
</template>

<style scoped>
/* Component styles are now handled by Tailwind CSS */
</style>
