import { useSubscription } from '@urql/vue'
import { computed, ref } from 'vue'

import { graphql } from '~/gql/gql.js'

export function useDotaState() {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)

	const {
		data,
		executeSubscription: connectSubscription,
		pause,
	} = useSubscription({
		query: graphql(`
			subscription DotaStateOverlay($apiKey: String!) {
				dotaState(apiKey: $apiKey) {
					channelId
					inGame
					mmr
					sessionWins
					sessionLosses
					winProbability
					heroName
					matchId
					teamIsRadiant
					teamKnown
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: paused,
	})

	const state = computed(() => data.value?.dotaState ?? null)

	function destroy() {
		pause()
	}

	function connect(key: string) {
		apiKey.value = key
		connectSubscription()
	}

	return {
		connect,
		destroy,
		state,
	}
}
