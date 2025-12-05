import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { ObsWebsocketModule } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type ObsWebsocketData = Omit<ObsWebsocketModule, '__typename'>

const invalidationKey = 'ObsWebsocketInvalidateKey'

export const useObsWebsocketApi = createGlobalState(() => {
	const useQueryObsWebsocket = () =>
		useQuery({
			variables: {},
			context: { additionalTypenames: [invalidationKey] },
			query: graphql(`
				query ObsWebsocketData {
					obsWebsocketData {
						serverPort
						serverAddress
						serverPassword
						sources
						audioSources
						scenes
						isConnected
					}
				}
			`),
		})

	const useMutationUpdateObsWebsocket = () =>
		useMutation(
			graphql(`
				mutation ObsWebsocketUpdate($input: ObsWebsocketUpdateInput!) {
					obsWebsocketUpdate(input: $input)
				}
			`),
			[invalidationKey]
		)

	const useSubscriptionIsConnected = () => {
		return useSubscription({
			query: graphql(`
				subscription ObsWebsocketIsConnected {
					obsWebsocketIsConnected
				}
			`),
			variables: {},
		})
	}

	return {
		useQueryObsWebsocket,
		useMutationUpdateObsWebsocket,
		useSubscriptionIsConnected,
	}
})
