import { useIntegrationsPageData } from '#layers/dashboard/api/integrations/integrations-page.ts'
import { useValorantIntegrationApi } from '#layers/dashboard/api/integrations/valorant.ts'
import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

export const valorantBroadcaster = new BroadcastChannel('valorant_channel')

export const useValorantIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()
	const { useLogoutMutation, usePostCodeMutation } = useValorantIntegrationApi()

	const { executeMutation: postCodeMutation } = usePostCodeMutation()
	const { executeMutation: logoutMutation } = useLogoutMutation()

	const isEnabled = computed(() => {
		return integrationsPage.valorantData.value?.enabled ?? false
	})

	const userName = computed(() => {
		return integrationsPage.valorantData.value?.userName ?? null
	})

	const avatar = computed(() => {
		return integrationsPage.valorantData.value?.avatar ?? null
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
		valorantBroadcaster.postMessage('refresh')
	}

	const authLink = computed(() => {
		return integrationsPage.valorantAuthLink.value
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
