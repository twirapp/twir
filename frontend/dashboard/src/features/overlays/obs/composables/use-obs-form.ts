import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, watchEffect } from 'vue'
import { z } from 'zod'

import { useObsWebsocketApi } from '@/api/overlays-obs'

const obsSettingsSchema = z.object({
	serverAddress: z.string().min(1, 'Server address is required').trim().default('localhost'),
	serverPort: z.number().min(1, 'Port is required').default(4455),
	serverPassword: z.string().min(1, 'Password is required').trim(),
})

export function useObsForm() {
	const api = useObsWebsocketApi()
	const { data: queryData, fetching } = api.useQueryObsWebsocket()
	const { data: isConnectedData } = api.useSubscriptionIsConnected()
	const updateMutation = api.useMutationUpdateObsWebsocket()

	const form = useForm({
		validationSchema: toTypedSchema(obsSettingsSchema),
		initialValues: {
			serverAddress: 'localhost',
			serverPort: 4455,
			serverPassword: '',
		},
	})

	const currentSettings = computed(() => queryData.value?.obsWebsocketData)

	// Use subscription for real-time connection status, fallback to query data
	const isConnected = computed(() => {
		// Prefer subscription data for real-time updates
		if (isConnectedData.value?.obsWebsocketIsConnected !== undefined) {
			return isConnectedData.value.obsWebsocketIsConnected
		}
		// Fallback to query data
		return currentSettings.value?.isConnected ?? false
	})

	// Use watchEffect to ensure it runs immediately and reactively
	watchEffect(() => {
		const newSettings = queryData.value?.obsWebsocketData
		if (newSettings && !fetching.value) {
			form.setFieldValue('serverAddress', newSettings.serverAddress)
			form.setFieldValue('serverPort', newSettings.serverPort)
			form.setFieldValue('serverPassword', newSettings.serverPassword)
		}
	})

	const onSubmit = form.handleSubmit(async (values) => {
		await updateMutation.executeMutation({
			input: {
				serverAddress: values.serverAddress,
				serverPort: values.serverPort,
				serverPassword: values.serverPassword,
			},
		})
	})

	return {
		form,
		onSubmit,
		isLoading: fetching,
		isSaving: computed(() => updateMutation.fetching.value),
		settings: currentSettings,
		isConnected,
	}
}
