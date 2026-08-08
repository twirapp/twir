import { createGlobalState } from '@vueuse/core'

import {
	integrationsPageCacheKey,
	useIntegrationsPageData,
} from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

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

	const donatexPostCode = () =>
		useMutation(
			graphql(`
				mutation DonateXPostCode($code: String!) {
					donatexPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const donatexLogout = () =>
		useMutation(
			graphql(`
				mutation DonateXLogout {
					donatexLogout
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

	const streamlabsPostCode = () =>
		useMutation(
			graphql(`
				mutation StreamlabsPostCode($code: String!) {
					streamlabsPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const streamlabsLogout = () =>
		useMutation(
			graphql(`
				mutation StreamlabsLogout {
					streamlabsLogout
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
		donatexPostCode,
		donatexLogout,
		vkPostCode,
		vkLogout,
		streamlabsPostCode,
		streamlabsLogout,

		broadcastRefresh,
	}
})
