import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import { graphql } from '@/gql'

export const integrationsPageCacheKey = 'integrationsPage'

/**
 * Unified query for all integrations page data.
 * This fetches all GraphQL-based integration data in a single request.
 */
const IntegrationsPageQuery = graphql(`
	query IntegrationsPageData {
		# Discord
		discordIntegrationData {
			guilds {
				id
				name
				icon
				liveNotificationEnabled
				liveNotificationChannelsIds
				liveNotificationShowTitle
				liveNotificationShowCategory
				liveNotificationShowViewers
				liveNotificationMessage
				liveNotificationShowPreview
				liveNotificationShowProfileImage
				offlineNotificationMessage
				shouldDeleteMessageOnOffline
				additionalUsersIdsForLiveCheck
			}
		}
		discordIntegrationAuthLink

		# Valorant
		valorantData {
			enabled
			userName
			avatar
		}
		valorantAuthLink

		# LastFM
		lastfmData {
			enabled
			userName
			avatar
		}
		lastfmAuthLink

		# Spotify
		spotifyData {
			userName
			avatar
		}
		spotifyAuthLink

		# DonationAlerts
		donationAlerts {
			enabled
			userName
			avatar
		}
		donationAlertsAuthLink

		# Donatello
		donatello {
			integrationId
		}

		# DonateStream
		integrationsDonateStream {
			integrationId
		}

		# DonatePay
		donatePayIntegration {
			apiKey
			enabled
		}

		# VK
		vk {
			enabled
			userName
			avatar
		}
		vkAuthLink

		# Faceit
		faceit {
			enabled
			userName
			avatar
			game
		}
		faceitAuthLink
	}
`)

export const useIntegrationsPageData = createGlobalState(() => {
	const refreshBroadcaster = new BroadcastChannel('integrations_page_broadcast_channel')

	const query = useQuery({
		query: IntegrationsPageQuery,
		context: {
			additionalTypenames: [integrationsPageCacheKey],
		},
		variables: {},
	})

	// Discord
	const discordGuilds = computed(() => query.data.value?.discordIntegrationData?.guilds ?? [])
	const discordAuthLink = computed(() => query.data.value?.discordIntegrationAuthLink ?? null)

	// Valorant
	const valorantData = computed(() => query.data.value?.valorantData ?? null)
	const valorantAuthLink = computed(() => query.data.value?.valorantAuthLink ?? null)

	// LastFM
	const lastfmData = computed(() => query.data.value?.lastfmData ?? null)
	const lastfmAuthLink = computed(() => query.data.value?.lastfmAuthLink ?? null)

	// Spotify
	const spotifyData = computed(() => query.data.value?.spotifyData ?? null)
	const spotifyAuthLink = computed(() => query.data.value?.spotifyAuthLink ?? null)

	// DonationAlerts
	const donationAlertsData = computed(() => query.data.value?.donationAlerts ?? null)
	const donationAlertsAuthLink = computed(() => query.data.value?.donationAlertsAuthLink ?? null)

	// Donatello
	const donatelloData = computed(() => query.data.value?.donatello ?? null)

	// DonateStream
	const donateStreamData = computed(() => query.data.value?.integrationsDonateStream ?? null)

	// DonatePay
	const donatePayData = computed(() => query.data.value?.donatePayIntegration ?? null)

	// VK
	const vkData = computed(() => query.data.value?.vk ?? null)
	const vkAuthLink = computed(() => query.data.value?.vkAuthLink ?? null)

	// Faceit
	const faceitData = computed(() => query.data.value?.faceit ?? null)
	const faceitAuthLink = computed(() => query.data.value?.faceitAuthLink ?? null)

	async function refetch() {
		await query.executeQuery({ requestPolicy: 'network-only' })
	}

	refreshBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		refetch()
	}

	function broadcastRefresh() {
		refreshBroadcaster.postMessage('refresh')
	}

	return {
		query,
		data: query.data,
		fetching: readonly(query.fetching),
		error: query.error,
		refetch,
		broadcastRefresh,

		// Discord
		discordGuilds,
		discordAuthLink,

		// Valorant
		valorantData,
		valorantAuthLink,

		// LastFM
		lastfmData,
		lastfmAuthLink,

		// Spotify
		spotifyData,
		spotifyAuthLink,

		// DonationAlerts
		donationAlertsData,
		donationAlertsAuthLink,

		// Donatello
		donatelloData,

		// DonateStream
		donateStreamData,

		// DonatePay
		donatePayData,

		// VK
		vkData,
		vkAuthLink,

		// Faceit
		faceitData,
		faceitAuthLink,
	}
})
