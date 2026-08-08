import { useSubscription } from '@urql/vue'
import { computed, ref } from 'vue'

import { graphql } from '~/gql/gql.js'

export function useStreamStatsGraphQL() {
	const apiKey = ref<string>('')
	const paused = computed(() => !apiKey.value)

	// Subscribe to settings updates
	const {
		data: settingsData,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription StreamStatsOverlaySettings($apiKey: String!) {
				overlaysStreamStats(apiKey: $apiKey) {
					id
					channelId
					design
					variant
					platformIconsEnabled
					viewersEnabled
					viewersMode
					viewersColor
					messagesEnabled
					messagesColor
					uptimeEnabled
					uptimeColor
					subscribersEnabled
					subscribersColor
					followersEnabled
					followersColor
					counterOrder
					customHtmlEnabled
					customHtml
					customCss
					createdAt
					updatedAt
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

	// Subscribe to live counters updates
	const {
		data: countersData,
		executeSubscription: connectCounters,
		pause: pauseCounters,
	} = useSubscription({
		query: graphql(`
			subscription StreamStatsOverlayCounters($apiKey: String!) {
				overlaysStreamStatsCounters(apiKey: $apiKey) {
					live
					viewers
					messages
					startedAt
					subscribers
					followers
					platformViewers {
						platform
						viewers
					}
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

	const settings = computed(() => settingsData.value?.overlaysStreamStats ?? null)
	const counters = computed(() => countersData.value?.overlaysStreamStatsCounters ?? null)

	function destroy() {
		pauseSettings()
		pauseCounters()
	}

	function connect(key: string) {
		apiKey.value = key
		connectSettings()
		connectCounters()
	}

	return {
		connect,
		destroy,
		settings,
		counters,
	}
}
