<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { type Font, FontSelector } from '@twir/fontsource'
import { InfoIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useBeRightBackOverlayApi, useProfile, useUserAccessFlagChecker } from '@/api'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { ColorPicker } from '@/components/ui/color-picker'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/toast'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import InputWithIcon from '@/components/ui/InputWithIcon.vue'
import { BeRightBackUpdateInputSchema } from '@/gql/validation-schemas.ts'

const { t } = useI18n()
const { toast } = useToast()

const { data: profile } = useProfile()

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

		toast({
			title: t('sharedTexts.saved'),
			variant: 'default',
		})
	} catch (e) {
		toast({
			title: 'Error occurred while saving BRB overlay',
			variant: 'destructive',
		})
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
		<Card class="flex-1">
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<h2 class="text-2xl font-bold">
					{{ t('overlays.brb.title', 'Be Right Back Overlay') }}
				</h2>
				<div class="flex gap-2">
					<Button variant="outline" @click="setDefaultSettings">
						{{ t('sharedButtons.setDefaultSettings') }}
					</Button>
					<Button
						variant="outline"
						:disabled="isSettingsError || isSettingsLoading || !canCopyLink"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</Button>
					<Button @click="save">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</CardHeader>

			<CardContent class="space-y-6">
				<form @submit="save" class="space-y-6">
					<div class="space-y-4">
						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.brb.settings.main.label') }}
						</h3>

						<div class="space-y-4">
							<div class="flex flex-col gap-2">
								<CommandButton
									name="brb"
									:title="t('overlays.brb.settings.main.startCommand.description')"
								/>
								<Alert>
									<InfoIcon class="h-4 w-4" />
									<AlertDescription>
										<span v-html="t('overlays.brb.settings.main.startCommand.example')" />
									</AlertDescription>
								</Alert>
								<CommandButton
									name="brbstop"
									:title="t('overlays.brb.settings.main.stopCommand.description')"
								/>
							</div>

							<FormField v-slot="{ componentField }" name="text">
								<FormItem>
									<FormLabel>{{ t('overlays.brb.settings.main.text') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" :maxlength="500" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ componentField }" name="backgroundColor">
								<FormItem>
									<FormLabel>{{ t('overlays.brb.settings.main.background') }}</FormLabel>
									<FormControl>
										<InputWithIcon
											id="backgroundColor"
											v-model="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										>
											<ColorPicker
												output-format="rgb"
												v-model="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</InputWithIcon>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ componentField }" name="fontColor">
								<FormItem>
									<FormLabel>{{ t('overlays.brb.settings.main.font.color') }}</FormLabel>
									<FormControl>
										<InputWithIcon
											id="fontColor"
											v-model="componentField.modelValue"
											@update:model-value="componentField['onUpdate:modelValue']"
										>
											<ColorPicker
												v-model="componentField.modelValue"
												@update:model-value="componentField['onUpdate:modelValue']"
											/>
										</InputWithIcon>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<div class="space-y-2">
								<Label>{{ t('overlays.brb.settings.main.font.family') }}</Label>
								<FontSelector
									v-model:font="fontData"
									:font-family="brbForm.values.fontFamily || 'inter'"
									font-style="normal"
									:font-weight="400"
								/>
							</div>

							<FormField v-slot="{ componentField }" name="fontSize">
								<FormItem>
									<FormLabel>{{ t('overlays.brb.settings.main.font.size') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" type="number" :min="1" :max="500" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>

						<Separator />
						<h3 class="text-lg font-semibold">
							{{ t('overlays.brb.settings.late.label') }}
						</h3>

						<div class="space-y-4">
							<FormField v-slot="{ componentField }" name="late.text">
								<FormItem>
									<FormLabel>{{ t('overlays.brb.settings.late.text') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" :maxlength="500" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ field }" name="late.enabled">
								<FormItem>
									<div class="flex items-center gap-2">
										<FormControl>
											<Switch
												:checked="field.value"
												@update:checked="field['onUpdate:modelValue']"
											/>
										</FormControl>
										<FormLabel class="!mt-0">{{ t('sharedTexts.enabled') }}</FormLabel>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ field }" name="late.displayBrbTime">
								<FormItem>
									<div class="flex items-center gap-2">
										<FormControl>
											<Switch
												:checked="field.value"
												@update:checked="field['onUpdate:modelValue']"
											/>
										</FormControl>
										<FormLabel class="!mt-0">{{
											t('overlays.brb.settings.late.displayBrb')
										}}</FormLabel>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</div>
				</form>
			</CardContent>
		</Card>

		<div class="flex-1 relative">
			<iframe
				v-if="brbIframeUrl"
				ref="brbIframeRef"
				:src="brbIframeUrl"
				class="w-full h-[600px] border rounded-lg"
			/>
			<div class="absolute top-4 right-4 flex gap-2">
				<Button variant="outline" size="sm" @click="sendIframeMessage('stop')">
					{{ t('overlays.brb.preview.stop') }}
				</Button>
				<Button
					size="sm"
					@click="
						() => {
							sendSettings()
							sendIframeMessage('start', { minutes: 0.1 })
						}
					"
				>
					{{ t('overlays.brb.preview.start') }}
				</Button>
			</div>
		</div>
	</div>
</template>

<style scoped>
/* Component styles are now handled by Tailwind CSS */
</style>
