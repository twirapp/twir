import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { ChatOverlaySettingsSubscription } from '@/gql/graphql.js'
import type { Settings } from '@twir/frontend-chat'

import { graphql } from '@/gql'

export const useChatOverlaySocket = createGlobalState(() => {
	const route = useRoute()
	const overlaySettings = ref<ChatOverlaySettingsSubscription['chatOverlaySettings'] | null>(null)

	const { data: neededData } = useQuery({
		query: graphql(`
			query ChatOverlayWithAdditionalData {
				authenticatedUser {
					id
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
				}
				twitchGetGlobalBadges {
					badges {
						set_id
						versions {
							id
							image_url_1x
							image_url_2x
							image_url_4x
						}
					}
				}
				twitchGetChannelBadges {
					badges {
						set_id
						versions {
							id
							image_url_1x
							image_url_2x
							image_url_4x
						}
					}
				}
			}
		`),
	})

	const sub = useSubscription({
		query: graphql(`
			subscription ChatOverlaySettings($id: String!, $apiKey: String!) {
				chatOverlaySettings(id: $id, apiKey: $apiKey) {
					id
					messageHideTimeout
					messageShowDelay
					preset
					fontSize
					hideCommands
					hideBots
					fontFamily
					showBadges
					showAnnounceBadge
					textShadowColor
					textShadowSize
					chatBackgroundColor
					direction
					fontWeight
					fontStyle
					paddingContainer
					animation
				}
			}
		`),
		variables: {
			id: route.query.id as string,
			apiKey: route.params.apiKey as string,
		},
	})

	watch(sub.data, (n) => {
		if (!n?.chatOverlaySettings) return

		overlaySettings.value = {
			...n.chatOverlaySettings,
		}
	})

	const chatLibSettings = computed<Settings | null>(() => {
		if (!overlaySettings.value || !neededData.value) return null

		return {
			...overlaySettings.value,
			channelBadges: neededData.value.twitchGetChannelBadges.badges,
			globalBadges: neededData.value.twitchGetGlobalBadges.badges,
			channelId: neededData.value.authenticatedUser.id ?? '',
			channelName: neededData.value.authenticatedUser.twitchProfile.login ?? '',
			channelDisplayName: neededData.value.authenticatedUser.twitchProfile.displayName ?? '',
		}
	})

	return {
		neededData,
		overlaySettings,
		chatLibSettings,
	}
})
