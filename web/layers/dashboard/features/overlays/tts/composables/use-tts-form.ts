import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, watchEffect } from 'vue'
import { z } from 'zod'

import { useTTSOverlayApi } from '#layers/dashboard/api/overlays-tts'

const ttsSettingsSchema = z.object({
	// General settings
	enabled: z.boolean(),
	voice: z.string().min(1, 'Voice is required'),
	pitch: z.number().min(0).max(100),
	rate: z.number().min(0).max(100),
	volume: z.number().min(0).max(100),
	// Advanced settings
	disallowedVoices: z.array(z.string()),
	doNotReadTwitchEmotes: z.boolean(),
	doNotReadEmoji: z.boolean(),
	doNotReadLinks: z.boolean(),
	allowUsersChooseVoiceInMainCommand: z.boolean(),
	maxSymbols: z.number().min(0).max(5000),
	readChatMessages: z.boolean(),
	readChatMessagesNicknames: z.boolean(),
})

export type TTSSettingsFormValues = z.infer<typeof ttsSettingsSchema>

export function useTTSForm() {
	const api = useTTSOverlayApi()
	const { data: settings, fetching } = api.useQueryTTS()
	const updateMutation = api.useMutationUpdateTTS()

	const form = useForm({
		validationSchema: toTypedSchema(ttsSettingsSchema),
		initialValues: {
			// General settings
			enabled: false,
			voice: 'alan',
			pitch: 50,
			rate: 50,
			volume: 30,
			// Advanced settings
			disallowedVoices: [],
			doNotReadTwitchEmotes: true,
			doNotReadEmoji: true,
			doNotReadLinks: true,
			allowUsersChooseVoiceInMainCommand: false,
			maxSymbols: 0,
			readChatMessages: false,
			readChatMessagesNicknames: false,
		},
	})

	const currentSettings = computed(() => settings.value?.overlaysTTS)

	// Use watchEffect to ensure it runs immediately and reactively
	watchEffect(() => {
		const newSettings = currentSettings.value
		if (newSettings && !fetching.value) {
			// General settings
			form.setFieldValue('enabled', newSettings.enabled)
			form.setFieldValue('voice', newSettings.voice)
			form.setFieldValue('pitch', newSettings.pitch)
			form.setFieldValue('rate', newSettings.rate)
			form.setFieldValue('volume', newSettings.volume)
			// Advanced settings
			form.setFieldValue('disallowedVoices', newSettings.disallowedVoices)
			form.setFieldValue('doNotReadTwitchEmotes', newSettings.doNotReadTwitchEmotes)
			form.setFieldValue('doNotReadEmoji', newSettings.doNotReadEmoji)
			form.setFieldValue('doNotReadLinks', newSettings.doNotReadLinks)
			form.setFieldValue('allowUsersChooseVoiceInMainCommand', newSettings.allowUsersChooseVoiceInMainCommand)
			form.setFieldValue('maxSymbols', newSettings.maxSymbols)
			form.setFieldValue('readChatMessages', newSettings.readChatMessages)
			form.setFieldValue('readChatMessagesNicknames', newSettings.readChatMessagesNicknames)
		}
	})

	const onSubmit = form.handleSubmit(async (values) => {
		await updateMutation.executeMutation({
			input: {
				enabled: values.enabled,
				voice: values.voice,
				pitch: values.pitch,
				rate: values.rate,
				volume: values.volume,
				disallowedVoices: values.disallowedVoices,
				doNotReadTwitchEmotes: values.doNotReadTwitchEmotes,
				doNotReadEmoji: values.doNotReadEmoji,
				doNotReadLinks: values.doNotReadLinks,
				allowUsersChooseVoiceInMainCommand: values.allowUsersChooseVoiceInMainCommand,
				maxSymbols: values.maxSymbols,
				readChatMessages: values.readChatMessages,
				readChatMessagesNicknames: values.readChatMessagesNicknames,
			},
		})
	})

	return {
		form,
		onSubmit,
		isLoading: fetching,
		isSaving: computed(() => updateMutation.fetching.value),
	}
}
