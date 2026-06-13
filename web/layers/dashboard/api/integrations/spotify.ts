import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'
import { integrationsPageCacheKey, useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'

export const useSpotifyIntegration = createGlobalState(() => {
	const spotifyBroadcaster = new BroadcastChannel('spotify_channel')
	const integrationsPage = useIntegrationsPageData()

	const postCode = useMutation(
		graphql(`
			mutation SpotifyPostCode($input: SpotifyPostCodeInput!) {
				spotifyPostCode(input: $input)
			}
		`),
		[integrationsPageCacheKey]
	)

	const logout = useMutation(
		graphql(`
			mutation SpotifyLogout {
				spotifyLogout
			}
		`),
		[integrationsPageCacheKey]
	)

	spotifyBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		integrationsPage.refetch()
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
