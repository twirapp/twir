import { createGlobalState } from '@vueuse/core'
import { useQuery } from '@urql/vue'

import { graphql } from '@/gql'
import { useMutation } from '@/composables/use-mutation.ts'

export const integrationsCacheKey = 'integrations'

export const useIntegrations = createGlobalState(() => {
	const refreshBroadcaster = new BroadcastChannel('integrations_broadcast_channel')

	const useData = (pause = false) =>
		useQuery({
			query: graphql(`
				query Integrations {
					donatello {
						integrationId
					}
					spotifyData {
						userName
						avatar
					}
					spotifyAuthLink
					integrationsDonateStream {
						integrationId
					}
					donationAlerts {
						enabled
						userName
						avatar
					}
					donationAlertsAuthLink
				}
			`),
			context: {
				additionalTypenames: [integrationsCacheKey],
			},
			variables: {},
			pause,
		})

	const donateStreamPostCode = () =>
		useMutation(
			graphql(`
				mutation DonateStreamPostCode($secret: String!) {
					integrationsDonateStreamPostSecret(input: { secret: $secret })
				}
			`)
		)

	const donationAlertsPostCode = () =>
		useMutation(
			graphql(`
				mutation DonationAlertsPostCode($code: String!) {
					donationAlertsPostCode(code: $code)
				}
			`),
			[integrationsCacheKey]
		)

	const donationAlertsLogout = () =>
		useMutation(
			graphql(`
				mutation DonationAlertsLogout {
					donationAlertsLogout
				}
			`),
			[integrationsCacheKey]
		)

	const { executeQuery: refreshIntegrations } = useData(true)

	refreshBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		refreshIntegrations({ requestPolicy: 'network-only' })
	}

	function broadcastRefresh() {
		refreshBroadcaster.postMessage('refresh')
	}

	return {
		useQuery: useData,
		donateStreamPostCode,
		donationAlertsPostCode,
		donationAlertsLogout,

		broadcastRefresh,
	}
})
