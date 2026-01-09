import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import { useIntegrationsPageData } from '#layers/dashboard/api/integrations/integrations-page.ts'
import { useLastfmIntegrationApi } from '#layers/dashboard/api/integrations/lastfm.ts'

export const lastfmBroadcaster = new BroadcastChannel('lastfm_channel')

export const useLastfmIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()
	const { useLogoutMutation, usePostCodeMutation } = useLastfmIntegrationApi()

	const { executeMutation: postCodeMutation } = usePostCodeMutation()
	const { executeMutation: logoutMutation } = useLogoutMutation()

	const isEnabled = computed(() => {
		return integrationsPage.lastfmData.value?.enabled ?? false
	})

	const userName = computed(() => {
		return integrationsPage.lastfmData.value?.userName ?? null
	})

	const avatar = computed(() => {
		return integrationsPage.lastfmData.value?.avatar ?? null
	})

	const isConfigured = computed(() => {
		return isEnabled.value && userName.value
	})

	async function refetchData() {
		await integrationsPage.refetch()
	}

	async function postCode(code: string) {
		const { error } = await postCodeMutation({ code })
		if (!error) {
			await refetchData()
		}
		return error
	}

	async function logout() {
		const { error } = await logoutMutation({})
		if (!error) {
			await refetchData()
		}
		return error
	}

	function broadcastRefresh() {
		lastfmBroadcaster.postMessage('refresh')
	}

	const authLink = computed(() => {
		return integrationsPage.lastfmAuthLink.value
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
		isDataFetching: readonly(integrationsPage.fetching),
	}
})
