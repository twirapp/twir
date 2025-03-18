import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'

const queryKey = ['spotify']

export const useSpotifyIntegration = createGlobalState(() => {
	const spotifyBroadcaster = new BroadcastChannel('spotify_channel')

	const spotifyData = useQuery({
		query: graphql(`
      query SpotifyData {
        profile: spotifyData {
          userName
          avatar
        }
				spotifyAuthLink
      }
    `),
		context: {
			additionalTypenames: queryKey,
		},
	})

	const postCode = useMutation(graphql(`
    mutation SpotifyPostCode($input: SpotifyPostCodeInput!) {
      spotifyPostCode(input: $input)
    }
  `), queryKey)

	const logout = useMutation(graphql(`
    mutation SpotifyLogout {
      spotifyLogout
    }
  `), queryKey)

	spotifyBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		spotifyData.executeQuery({ requestPolicy: 'network-only' })
	}

	function broadcastRefresh() {
		spotifyBroadcaster.postMessage('refresh')
	}

	return {
		spotifyData,
		postCode,
		logout,
		broadcastRefresh,
	}
})
