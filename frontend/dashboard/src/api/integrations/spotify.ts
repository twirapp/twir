import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'
import { integrationsCacheKey, useIntegrations } from '@/api/integrations/integrations.ts'

export const useSpotifyIntegration = createGlobalState(() => {
	const spotifyBroadcaster = new BroadcastChannel('spotify_channel')
	const integrationsManager = useIntegrations()
	const { executeQuery: refreshIntegrations } = integrationsManager.useQuery()

	const postCode = useMutation(
		graphql(`
			mutation SpotifyPostCode($input: SpotifyPostCodeInput!) {
				spotifyPostCode(input: $input)
			}
		`),
		[integrationsCacheKey]
	)

	const logout = useMutation(
		graphql(`
			mutation SpotifyLogout {
				spotifyLogout
			}
		`),
		[integrationsCacheKey]
	)

	spotifyBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		refreshIntegrations({ requestPolicy: 'network-only' })
	}

	function broadcastRefresh() {
		spotifyBroadcaster.postMessage('refresh')
	}

	return {
		postCode,
		logout,
		broadcastRefresh,
	}
})
