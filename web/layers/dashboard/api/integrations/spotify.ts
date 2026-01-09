import { createGlobalState } from '@vueuse/core'

import {
	integrationsPageCacheKey,
	useIntegrationsPageData,
} from '@/api/integrations/integrations-page.ts'
import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation.ts'

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
