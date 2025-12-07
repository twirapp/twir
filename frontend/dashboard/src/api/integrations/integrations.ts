import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'
import { useMutation } from '@/composables/use-mutation.ts'
import {
	integrationsPageCacheKey,
	useIntegrationsPageData,
} from '@/api/integrations/integrations-page.ts'

export const useIntegrations = createGlobalState(() => {
	const refreshBroadcaster = new BroadcastChannel('integrations_broadcast_channel')
	const integrationsPage = useIntegrationsPageData()

	const donateStreamPostCode = () =>
		useMutation(
			graphql(`
				mutation DonateStreamPostCode($secret: String!) {
					integrationsDonateStreamPostSecret(input: { secret: $secret })
				}
			`),
			[integrationsPageCacheKey]
		)

	const donationAlertsPostCode = () =>
		useMutation(
			graphql(`
				mutation DonationAlertsPostCode($code: String!) {
					donationAlertsPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const donationAlertsLogout = () =>
		useMutation(
			graphql(`
				mutation DonationAlertsLogout {
					donationAlertsLogout
				}
			`),
			[integrationsPageCacheKey]
		)

	const vkPostCode = () =>
		useMutation(
			graphql(`
				mutation VkPostCode($code: String!) {
					vkPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const vkLogout = () =>
		useMutation(
			graphql(`
				mutation VkLogout {
					vkLogout
				}
			`),
			[integrationsPageCacheKey]
		)

	refreshBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		integrationsPage.refetch()
	}

	function broadcastRefresh() {
		refreshBroadcaster.postMessage('refresh')
	}

	return {
		donateStreamPostCode,
		donationAlertsPostCode,
		donationAlertsLogout,
		vkPostCode,
		vkLogout,

		broadcastRefresh,
	}
})
