import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import { useLastfmIntegrationApi } from '@/api/integrations/lastfm.ts'

export const lastfmBroadcaster = new BroadcastChannel('lastfm_channel')

export const useLastfmIntegration = createGlobalState(() => {
	const { useData, useLogoutMutation, usePostCodeMutation } = useLastfmIntegrationApi()

	const { data: integrationData, executeQuery, fetching: isDataFetching } = useData()

	const { executeMutation: postCodeMutation } = usePostCodeMutation()
	const { executeMutation: logoutMutation } = useLogoutMutation()

	const isEnabled = computed(() => {
		return integrationData.value?.lastfmData?.enabled ?? false
	})

	const userName = computed(() => {
		return integrationData.value?.lastfmData?.userName ?? null
	})

	const avatar = computed(() => {
		return integrationData.value?.lastfmData?.avatar ?? null
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
		lastfmBroadcaster.postMessage('refresh')
	}

	const authLink = computed(() => {
		return integrationData.value?.lastfmAuthLink
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
