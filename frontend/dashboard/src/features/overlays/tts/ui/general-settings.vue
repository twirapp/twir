<script setup lang="ts">
import { InfoIcon, Volume2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useTTSForm } from '@/features/overlays/tts/composables/use-tts-form'
import { useTTSVoices } from '@/features/overlays/tts/composables/use-tts-voices'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Slider } from '@/components/ui/slider'
import { Switch } from '@/components/ui/switch'
import { toast } from 'vue-sonner'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink'
import { useProfile, useUserAccessFlagChecker } from '@/api'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const { form, onSubmit, isLoading, isSaving } = useTTSForm()
const { voices } = useTTSVoices()

const { data: profile } = useProfile()
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
		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<h2 class="text-2xl font-bold">
					{{ t('overlays.tts.generalSettings') }}
				</h2>
				<div class="flex gap-2">
					<Button
						variant="outline"
						:disabled="isLoading || !canCopyLink"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</Button>
					<Button :disabled="isLoading || isSaving" @click="handleSave">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</CardHeader>

			<CardContent class="space-y-6">
				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertDescription>
						{{ t('overlays.tts.eventsHint') }}
					</AlertDescription>
				</Alert>

				<form @submit.prevent="handleSave" class="space-y-6">
					<!-- Enabled Switch -->
					<FormField v-slot="{ field }" name="enabled">
						<FormItem>
							<div class="flex items-center justify-between">
								<FormLabel>{{ t('sharedTexts.enabled') }}</FormLabel>
								<FormControl>
									<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
								</FormControl>
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<Separator />

					<!-- Voice Selection -->
					<FormField v-slot="{ componentField }" name="voice">
						<FormItem>
							<FormLabel>{{ t('overlays.tts.voice') }}</FormLabel>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger>
										<SelectValue :placeholder="t('overlays.tts.selectVoice')" />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectItem v-for="voice in voices" :key="voice.value" :value="voice.value">
										{{ voice.label }}
									</SelectItem>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Voice Preview -->
					<div class="space-y-2">
						<Label>{{ t('overlays.tts.preview') }}</Label>
						<div class="flex gap-2">
							<Input
								v-model="previewText"
								:placeholder="t('overlays.tts.previewText')"
								class="flex-1"
							/>
							<Button
								type="button"
								variant="outline"
								size="icon"
								:disabled="!previewText || isPlaying"
								@click="playPreview"
							>
								<Volume2 class="h-4 w-4" />
							</Button>
						</div>
					</div>

					<Separator />

					<!-- Volume Slider -->
					<FormField v-slot="{ componentField }" name="volume">
						<FormItem>
							<div class="space-y-2">
								<FormLabel>{{ t('overlays.tts.volume') }}</FormLabel>
								<div class="flex items-center gap-4">
									<FormControl class="flex-1">
										<Slider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</FormControl>
									<Input
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
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Pitch Slider -->
					<FormField v-slot="{ componentField }" name="pitch">
						<FormItem>
							<div class="space-y-2">
								<FormLabel>{{ t('overlays.tts.pitch') }}</FormLabel>
								<div class="flex items-center gap-4">
									<FormControl class="flex-1">
										<Slider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</FormControl>
									<Input
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
							<FormMessage />
						</FormItem>
					</FormField>

					<!-- Rate Slider -->
					<FormField v-slot="{ componentField }" name="rate">
						<FormItem>
							<div class="space-y-2">
								<FormLabel>{{ t('overlays.tts.rate') }}</FormLabel>
								<div class="flex items-center gap-4">
									<FormControl class="flex-1">
										<Slider
											:model-value="[componentField.modelValue]"
											:min="0"
											:max="100"
											:step="1"
											@update:model-value="
												(val) =>
													val?.[0] !== undefined && componentField['onUpdate:modelValue']?.(val[0])
											"
										/>
									</FormControl>
									<Input
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
							<FormMessage />
						</FormItem>
					</FormField>
				</form>
			</CardContent>
		</Card>

		<!-- Advanced Settings Card -->
		<Card>
			<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-4">
				<div>
					<CardTitle>{{ t('overlays.tts.advancedSettings') }}</CardTitle>
					<CardDescription>{{ t('overlays.tts.advancedDescription') }}</CardDescription>
				</div>
				<Button :disabled="isLoading || isSaving" @click="handleSave">
					{{ t('sharedButtons.save') }}
				</Button>
			</CardHeader>

			<CardContent class="space-y-6">
				<form @submit.prevent="handleSave" class="space-y-6">
					<!-- Reading Options Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.readingOptions') }}
						</h3>

						<FormField v-slot="{ field }" name="doNotReadTwitchEmotes">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.doNotReadTwitchEmotes') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.doNotReadTwitchEmotesDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ field }" name="doNotReadEmoji">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.doNotReadEmoji') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.doNotReadEmojiDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ field }" name="doNotReadLinks">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.doNotReadLinks') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.doNotReadLinksDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<Separator />

					<!-- Chat Messages Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.chatMessages') }}
						</h3>

						<FormField v-slot="{ field }" name="readChatMessages">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.readChatMessages') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.readChatMessagesDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ field }" name="readChatMessagesNicknames">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.readChatMessagesNicknames') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.readChatMessagesNicknamesDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<Separator />

					<!-- User Settings Section -->
					<div class="space-y-4">
						<h3 class="text-lg font-semibold">
							{{ t('overlays.tts.userSettings') }}
						</h3>

						<FormField v-slot="{ field }" name="allowUsersChooseVoiceInMainCommand">
							<FormItem>
								<div class="flex items-center justify-between">
									<div class="space-y-1">
										<FormLabel>{{ t('overlays.tts.allowUsersChooseVoice') }}</FormLabel>
										<FormDescription>
											{{ t('overlays.tts.allowUsersChooseVoiceDesc') }}
										</FormDescription>
									</div>
									<FormControl>
										<Switch :checked="field.value" @update:checked="field['onUpdate:modelValue']" />
									</FormControl>
								</div>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField v-slot="{ componentField }" name="maxSymbols">
							<FormItem>
								<FormLabel>{{ t('overlays.tts.maxSymbols') }}</FormLabel>
								<FormDescription>
									{{ t('overlays.tts.maxSymbolsDesc') }}
								</FormDescription>
								<FormControl>
									<Input v-bind="componentField" type="number" :min="0" :max="5000" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>
				</form>
			</CardContent>
		</Card>
	</div>
</template>
