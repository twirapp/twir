import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import type { ObsWebsocketModule } from '@/gql/graphql.js'

import { useProfile } from '@/api/auth.js'
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

	// Subscribe to OBS data updates using apiKey from selected dashboard
	const useSubscriptionObsData = () => {
		const { data: profile } = useProfile()

		const apiKey = computed(() => {
			const selectedDashboard = profile.value?.availableDashboards.find(
				(d) => d.id === profile.value?.selectedDashboardId
			)
			return selectedDashboard?.apiKey ?? ''
		})

		const paused = computed(() => !apiKey.value)

		return useSubscription({
			query: graphql(`
				subscription ObsWebsocketDataDashboard($apiKey: String!) {
					obsWebsocketData(apiKey: $apiKey) {
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
			get variables() {
				return { apiKey: apiKey.value }
			},
			pause: paused,
		})
	}

	return {
		useQueryObsWebsocket,
		useMutationUpdateObsWebsocket,
		useSubscriptionIsConnected,
		useSubscriptionObsData,
	}
})
