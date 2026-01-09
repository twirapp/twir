<script setup lang="ts">
import { InfoIcon, Volume2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'


import { useTTSForm } from '~/features/overlays/tts/composables/use-tts-form'
import { useTTSVoices } from '~/features/overlays/tts/composables/use-tts-voices'










import { toast } from 'vue-sonner'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink'
import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
const { form, onSubmit, isLoading, isSaving } = useTTSForm()
const { voices } = useTTSVoices()

const { user: profile } = storeToRefs(useDashboardAuth())
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const { copyOverlayLink } = useCopyOverlayLink('tts')

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays.value
})

const previewText = ref('Hello world, this is a test message')
const isPlaying = ref(false)

async function handleSave() {
	await onSubmit()
	toast.success(t('sharedTexts.saved'))
}

async function playPreview() {
	if (!previewText.value || isPlaying.value) return

	isPlaying.value = true

	try {
		const formValues = form.values as any
		const voice = formValues.voice || 'alan'
		const pitch = formValues.pitch || 50
		const rate = formValues.rate || 50
		const volume = formValues.volume || 50

		// Build query string
		const params = new URLSearchParams({
			voice,
			text: previewText.value,
			pitch: pitch.toString(),
			rate: rate.toString(),
			volume: volume.toString(),
		})

		// Fetch audio directly as blob since the generated API client doesn't handle binary responses well
		const response = await fetch(`/api/v1/tts/say?${params.toString()}`, {
			credentials: 'include',
		})

		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`)
		}

		const blob = await response.blob()
		const audioUrl = URL.createObjectURL(blob)
		const audio = new Audio(audioUrl)

		audio.onended = () => {
			isPlaying.value = false
			URL.revokeObjectURL(audioUrl)
		}

		audio.onerror = () => {
			isPlaying.value = false
			URL.revokeObjectURL(audioUrl)
			toast({
				title: t('sharedTexts.error'),
				description: t('overlays.tts.previewError'),
				variant: 'destructive',
			})
		}

		await audio.play()
	} catch (error) {
		console.error('Failed to play TTS preview:', error)
		isPlaying.value = false
		toast({
			title: t('sharedTexts.error'),
			description: t('overlays.tts.previewError'),
			variant: 'destructive',
		})
	}
}
</script>

<template>
	<div class="flex flex-col gap-4 p-4">
		<!-- General Settings Card -->
		<UiCard>
			<UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<h2 class="text-2xl font-bold">
					{{ t('overlays.tts.generalSettings') }}
				</h2>
				<div class="flex gap-2">
					<UiButton
						variant="outline"
						:disabled="isLoading || !canCopyLink"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</UiButton>
					<UiButton :disabled="isLoading || isSaving" @click="handleSave">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</div>
			</UiCardHeader>

			<UiCardContent class="space-y-6">
				<UiAlert>
					<InfoIcon class="h-4 w-4" />
					<UiAlertDescription>
						{{ t('overlays.tts.eventsHint') }}
					</UiAlertDescription>
				</UiAlert>

				<form @submit.prevent="handleSave" class="space-y-6">
					<!-- Enabled Switch -->
					<UiFormField v-slot="{ field }" name="enabled">
						<UiFormItem>
							<div class="flex items-center justify-between">
								<UiFormLabel>{{ t('sharedTexts.enabled') }}</UiFormLabel>
								<UiFormControl>
									<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
								</UiFormControl>
							</div>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiSeparator />

					<!-- Voice Selection -->
					<UiFormField v-slot="{ componentField }" name="voice">
						<UiFormItem>
							<UiFormLabel>{{ t('overlays.tts.voice') }}</UiFormLabel>
							<UiSelect v-bind="componentField">
								<UiFormControl>
									<UiSelectTrigger>
										<UiSelectValue :placeholder="t('overlays.tts.selectVoice')" />
									</UiSelectTrigger>
								</UiFormControl>
								<UiSelectContent>
									<UiSelectItem v-for="voice in voices" :key="voice.value" :value="voice.value">
										{{ voice.label }}
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<!-- Voice Preview -->
					<div class="space-y-2">
						<UiLabel>{{ t('overlays.tts.preview') }}</UiLabel>
						<div class="flex gap-2">
							<UiInput
								v-model="previewText"
								:placeholder="t('overlays.tts.previewText')"
								class="flex-1"
							/>
							<UiButton
								type="button"
								variant="outline"
								size="icon"
								:disabled="!previewText || isPlaying"
								@click="playPreview"
							>
								<Volume2 class="h-4 w-4" />
							</UiButton>
						</div>
					</div>

					<UiSeparator />

					<!-- Volume Slider -->
					<UiFormField v-slot="{ componentField }" name="volume">
						<UiFormItem>
							<div class="space-y-2">
								<UiFormLabel>{{ t('overlays.tts.volume') }}</UiFormLabel>
								<div class="flex items-center gap-4">
									<UiFormControl class="flex-1">
										<UiSlider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</UiFormControl>
									<UiInput
										:model-value="componentField.modelValue"
										@update:model-value="
											(val) => componentField['onUpdate:modelValue']?.(Number(val))
										"
										type="number"
										:min="0"
										:max="100"
										class="w-20"
									/>
								</div>
							</div>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<!-- Pitch Slider -->
					<UiFormField v-slot="{ componentField }" name="pitch">
						<UiFormItem>
							<div class="space-y-2">
								<UiFormLabel>{{ t('overlays.tts.pitch') }}</UiFormLabel>
								<div class="flex items-center gap-4">
									<UiFormControl class="flex-1">
										<UiSlider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</UiFormControl>
									<UiInput
										:model-value="componentField.modelValue"
										@update:model-value="
											(val) => componentField['onUpdate:modelValue']?.(Number(val))
										"
										type="number"
										:min="0"
										:max="100"
										class="w-20"
									/>
								</div>
							</div>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<!-- Rate Slider -->
					<UiFormField v-slot="{ componentField }" name="rate">
						<UiFormItem>
							<div class="space-y-2">
								<UiFormLabel>{{ t('overlays.tts.rate') }}</UiFormLabel>
								<div class="flex items-center gap-4">
									<UiFormControl class="flex-1">
										<UiSlider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</UiFormControl>
									<UiInput
										:model-value="componentField.modelValue"
										@update:model-value="
											(val) => componentField['onUpdate:modelValue']?.(Number(val))
										"
										type="number"
										:min="0"
										:max="100"
										class="w-20"
									/>
								</div>
							</div>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</form>
			</UiCardContent>
		</UiCard>

		<!-- Advanced Settings Card -->
		<UiCard>
			<UiCardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<div>
					<UiCardTitle>{{ t('overlays.tts.advancedSettings') }}</UiCardTitle>
					<UiCardDescription>{{ t('overlays.tts.advancedDescription') }}</UiCardDescription>
				</div>
				<UiButton :disabled="isLoading || isSaving" @click="handleSave">
					{{ t('sharedButtons.save') }}
				</UiButton>
			</UiCardHeader>

			<UiCardContent class="space-y-6">
				<form @submit.prevent="handleSave" class="space-y-6">
					<!-- Reading Options Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.readingOptions') }}
						</h3>

						<UiFormField v-slot="{ field }" name="doNotReadTwitchEmotes">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.doNotReadTwitchEmotes') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.doNotReadTwitchEmotesDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ field }" name="doNotReadEmoji">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.doNotReadEmoji') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.doNotReadEmojiDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ field }" name="doNotReadLinks">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.doNotReadLinks') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.doNotReadLinksDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</div>

					<UiSeparator />

					<!-- Chat Messages Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.chatMessages') }}
						</h3>

						<UiFormField v-slot="{ field }" name="readChatMessages">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.readChatMessages') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.readChatMessagesDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ field }" name="readChatMessagesNicknames">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.readChatMessagesNicknames') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.readChatMessagesNicknamesDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</div>

					<UiSeparator />

					<!-- User Settings Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.userSettings') }}
						</h3>

						<UiFormField v-slot="{ field }" name="allowUsersChooseVoiceInMainCommand">
							<UiFormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<UiFormLabel>{{ t('overlays.tts.allowUsersChooseVoice') }}</UiFormLabel>
										<UiFormDescription>
											{{ t('overlays.tts.allowUsersChooseVoiceDesc') }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch :model-value="field.value" @update:model-value="field['onUpdate:modelValue']" />
									</UiFormControl>
								</div>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="maxSymbols">
							<UiFormItem>
								<UiFormLabel>{{ t('overlays.tts.maxSymbols') }}</UiFormLabel>
								<UiFormDescription>
									{{ t('overlays.tts.maxSymbolsDesc') }}
								</UiFormDescription>
								<UiFormControl>
									<UiInput v-bind="componentField" type="number" :min="0" :max="5000" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</div>
				</form>
			</UiCardContent>
		</UiCard>
	</div>
</template>
