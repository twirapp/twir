import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState, useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { TwirWebSocketEvent } from '@/api.js'
import type { BadgeVersion, ChatBadge, Settings } from '@twir/frontend-chat'

import { graphql } from '@/gql'
import { generateSocketUrlWithParams } from '@/helpers.js'

export const useChatOverlaySocket = createGlobalState(() => {
	const route = useRoute()

	useQuery({
		query: graphql(`
			query ChatOverlayWithAdditionalData($overlayId: String!) {
				authenticatedUser {
					id
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
				}
				chatOverlaysById(id: $overlayId) {
					id
					messageHideTimeout
				}
				twitchGetGlobalBadges {
					badges {
						set_id
						versions {
							image_url_1x
						}
					}
				}
				twitchGetChannelBadges {
					badges {
						set_id
						versions {
							image_url_1x
						}
					}
				}
			}
		`),
		variables: {
			overlayId: route.query.id,
		},
	})

	const sub = useSubscription({
		query: graphql(`
			subscription ChatOverlaySettings($id: String!, $apiKey: String!) {
				chatOverlaySettings(id: $id, apiKey: $apiKey) {
					id
					animation
					fontSize
				}
			}
		`),
		variables: {
			id: route.query.id,
			apiKey: route.params.apiKey,
		},
	})

	watch(sub.data, (n) => console.log(n))

	const settings = ref<Settings>({
		channelId: '',
		channelName: '',
		channelDisplayName: '',
		globalBadges: new Map<string, ChatBadge>(),
		channelBadges: new Map<string, BadgeVersion>(),
		messageHideTimeout: 0,
		messageShowDelay: 0,
		preset: 'clean',
		fontSize: 20,
		hideBots: false,
		hideCommands: false,
		fontFamily: 'Roboto',
		showAnnounceBadge: true,
		showBadges: true,
		textShadowColor: '',
		textShadowSize: 0,
		chatBackgroundColor: '',
		direction: 'top',
		fontStyle: 'normal',
		fontWeight: 400,
		paddingContainer: 0,
	})

	const overlayId = ref<string | undefined>()
	const socketUrl = ref('')

	const { data, status, send, open, close } = useWebSocket(
		socketUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }))
			},
		},
	)

	watch(data, (d: string) => {
		const event = JSON.parse(d) as TwirWebSocketEvent<Settings>
		if (event.eventName === 'settings') {
			if (overlayId.value && event.data.id !== overlayId.value) return

			const data = event.data

			settings.value = {
				...data,
				globalBadges: new Map(),
				channelBadges: new Map(),
			}

			for (const badge of Object.values(data.globalBadges)) {
				settings.value.globalBadges.set(badge.set_id, badge)
			}

			for (const [setId, version] of Object.entries(data.channelBadges)) {
				settings.value.channelBadges.set(setId, version)
			}
		}
	})

	function destroy(): void {
		close()
	}

	function connect(apiKey: string, _overlayId?: string): void {
		if (status.value === 'OPEN') return

		const url = generateSocketUrlWithParams('/overlays/chat', {
			apiKey,
			id: _overlayId,
		})

		socketUrl.value = url
		overlayId.value = _overlayId

		open()
	}

	return {
		settings,
		connect,
		destroy,
	}
})
