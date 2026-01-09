import { useIntegrationsPageData } from '#layers/dashboard/api/integrations/integrations-page.ts'
import { useVKIntegrationApi } from '#layers/dashboard/api/integrations/vk.ts'
import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

export const vkBroadcaster = new BroadcastChannel('vk_channel')

export const useVKIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()
	const { useLogoutMutation, usePostCodeMutation } = useVKIntegrationApi()

	const { executeMutation: postCodeMutation } = usePostCodeMutation()
	const { executeMutation: logoutMutation } = useLogoutMutation()

	const isEnabled = computed(() => {
		return integrationsPage.vkData.value?.enabled ?? false
	})

	const userName = computed(() => {
		return integrationsPage.vkData.value?.userName ?? null
	})

	const avatar = computed(() => {
		return integrationsPage.vkData.value?.avatar ?? null
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
		vkBroadcaster.postMessage('refresh')
	}

	const authLink = computed(() => {
		return integrationsPage.vkAuthLink.value
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
