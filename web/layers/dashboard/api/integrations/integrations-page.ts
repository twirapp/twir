import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, readonly } from 'vue'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth.js'
import { graphql } from '~/gql/gql.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

export const integrationsPageCacheKey = 'integrationsPage'

/**
 * Unified query for all integrations page data.
 * This fetches all GraphQL-based integration data in a single request.
 */
export const IntegrationsPageQuery = graphql(`
	query IntegrationsPageData($canManageIntegrations: Boolean!) {
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
		discordIntegrationAuthLink @include(if: $canManageIntegrations)

		# Valorant
		valorantData {
			enabled
			userName
			avatar
		}
		valorantAuthLink @include(if: $canManageIntegrations)

		# LastFM
		lastfmData {
			enabled
			userName
			avatar
		}
		lastfmAuthLink @include(if: $canManageIntegrations)

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
		donationAlertsAuthLink @include(if: $canManageIntegrations)

		# Donatello
		donatello @include(if: $canManageIntegrations) {
			integrationId
		}

		# DonateStream
		integrationsDonateStream @include(if: $canManageIntegrations) {
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
		vkAuthLink @include(if: $canManageIntegrations)

		# Faceit
		faceit {
			enabled
			userName
			avatar
			game
			faceitUserId
		}
		faceitAuthLink @include(if: $canManageIntegrations)

		# Streamlabs
		streamlabs {
			enabled
			userName
			avatar
		}
		streamlabsAuthLink @include(if: $canManageIntegrations)

		# Imports
		nightbotGetData {
			userName
			avatar
		}
		nightbotGetAuthLink @include(if: $canManageIntegrations)
		streamelementsGetData {
			userName
			avatar
		}
		streamelementsGetAuthorizationUrl @include(if: $canManageIntegrations)
	}
`)

export const useIntegrationsPageData = createGlobalState(() => {
	const refreshBroadcaster = new BroadcastChannel('integrations_page_broadcast_channel')
	const canManageIntegrations = useUserAccessFlagChecker(
		ChannelRolePermissionEnum.ManageIntegrations
	)

	const query = useQuery({
		query: IntegrationsPageQuery,
		context: {
			additionalTypenames: [integrationsPageCacheKey],
		},
		variables: computed(() => ({ canManageIntegrations: canManageIntegrations.value })),
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

	// Streamlabs
	const streamlabsData = computed(() => query.data.value?.streamlabs ?? null)
	const streamlabsAuthLink = computed(() => query.data.value?.streamlabsAuthLink ?? null)

	// Imports
	const nightbotData = computed(() => query.data.value?.nightbotGetData ?? null)
	const nightbotAuthLink = computed(() => query.data.value?.nightbotGetAuthLink ?? null)
	const streamelementsData = computed(() => query.data.value?.streamelementsGetData ?? null)
	const streamelementsAuthLink = computed(
		() => query.data.value?.streamelementsGetAuthorizationUrl ?? null
	)

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

		// Streamlabs
		streamlabsData,
		streamlabsAuthLink,

		// Imports
		nightbotData,
		nightbotAuthLink,
		streamelementsData,
		streamelementsAuthLink,
	}
})
