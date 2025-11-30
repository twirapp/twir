import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import { useValorantIntegrationApi } from '@/api/integrations/valorant.ts'

export const valorantBroadcaster = new BroadcastChannel('valorant_channel')

export const useValorantIntegration = createGlobalState(() => {
	const { useData, useLogoutMutation, usePostCodeMutation } = useValorantIntegrationApi()

	const { data: integrationData, executeQuery, fetching: isDataFetching } = useData()

	const { executeMutation: postCodeMutation } = usePostCodeMutation()
	const { executeMutation: logoutMutation } = useLogoutMutation()

	const isEnabled = computed(() => {
		return integrationData.value?.valorantData?.enabled ?? false
	})

	const userName = computed(() => {
		return integrationData.value?.valorantData?.userName ?? null
	})

	const avatar = computed(() => {
		return integrationData.value?.valorantData?.avatar ?? null
	})

	const isConfigured = computed(() => {
		return isEnabled.value && userName.value
	})

	async function refetchData() {
		await executeQuery({ requestPolicy: 'network-only' })
	}

	async function postCode(code: string) {
		const { error } = await postCodeMutation({ code })
		return error
	}

	async function logout() {
		const { error } = await logoutMutation({})

		return error
	}

	function broadcastRefresh() {
		valorantBroadcaster.postMessage('refresh')
	}

	const authLink = computed(() => {
		return integrationData.value?.valorantAuthLink
	})

	return {
		isEnabled,
		userName,
		avatar,
		isConfigured,
		authLink,
		logout,
		postCode,
		broadcastRefresh,
		refetchData,
		isDataFetching: readonly(isDataFetching),
	}
})
