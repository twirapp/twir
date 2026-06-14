import { useSubscription } from '@urql/vue'
import { type MaybeRef, computed, ref, unref, watch } from 'vue'

import type { Buidler } from './use-kappagen-builder'
import type { KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { type ChatMessage, type ChatSettings, useChatTmi } from '~~/layers/overlays/composables/tmi/use-chat-tmi'
import { useMessageHelpers } from '~~/layers/overlays/composables/tmi/use-message-helpers'
import { useSevenTv } from '~~/layers/overlays/composables/tmi/use-seven-tv'
import { graphql } from '~~/app/gql/gql'

export function useKappagenOverlaySocket(
	instance: MaybeRef<KappagenMethods>,
	emotesBuilder: Buidler
) {
	const apiKey = ref<string>('')

	const paused = computed(() => !apiKey.value)

	const {
		data: eventsData,
		executeSubscription: connectEvents,
		pause: pauseEvents,
	} = useSubscription({
		query: graphql(`
			subscription TwirEvents($apiKey: String!) {
				twirEvents(apiKey: $apiKey) {
					baseInfo {
						channelId
						channelName
						type
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
	const {
		data: settings,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription KappagenSettings($apiKey: String!) {
				overlaysKappagen(apiKey: $apiKey) {
					id
					enableSpawn
					excludedEmotes
					enableRave
					animation {
						fadeIn
						fadeOut
						zoomIn
						zoomOut
					}
					animations {
						style
						prefs {
							size
							center
							speed
							faces
							message
							time
						}
						count
						enabled
					}
					emotes {
						time
						max
						queue
						ffzEnabled
						bttvEnabled
						sevenTvEnabled
						emojiStyle
					}
					size {
						rationNormal
						rationSmall
						min
						max
					}
					events {
						event
						disabledAnimations
						enabled
					}
					createdAt
					updatedAt
					channel {
						id
						login
					}
					channelIdentities {
						platform
						id
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
	const { data: triggerKappagenData } = useSubscription({
		query: graphql(`
			subscription KappagenTrigger($apiKey: String!) {
				overlaysKappagenTrigger(apiKey: $apiKey) {
					text
					emotes {
						id
						url
						positions
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
	const {
		data: chatMessagesData,
		executeSubscription: connectChatMessages,
		pause: pauseChatMessages,
	} = useSubscription({
		query: graphql(`
			subscription KappagenChatMessages($apiKey: String!) {
				overlaysKappagenChatMessages(apiKey: $apiKey) {
					id
					platform
					text
					channelId
					channelLogin
					channelName
					userID
					userName
					userDisplayName
					userColor
					createdAt
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

	const { makeMessageChunks } = useMessageHelpers()
	const { fetchSevenTvEmotes } = useSevenTv()

	function handleSpawnChunks(chunks: ReturnType<typeof makeMessageChunks>) {
		if (!settings.value?.overlaysKappagen.enableSpawn) return

		const firstChunk = chunks.at(0)
		if (firstChunk?.type === 'text' && firstChunk.value.startsWith('!')) {
			return
		}

		const generatedEmotes = emotesBuilder.buildSpawnEmotes(chunks)
		if (!generatedEmotes.length) return

		unref(instance).showEmotes(generatedEmotes)
	}

	function onTwitchMessage(msg: ChatMessage) {
		if (msg.type === 'system') return
		handleSpawnChunks(msg.chunks)
	}

	const twitchChannelIdentity = computed(() => {
		return settings.value?.overlaysKappagen.channelIdentities.find((identity) => identity.platform === 'TWITCH')
	})
	const kickChannelIdentity = computed(() => {
		return settings.value?.overlaysKappagen.channelIdentities.find((identity) => identity.platform === 'KICK')
	})

	const twitchChatSettings = computed<ChatSettings>(() => {
		return {
			channelId: twitchChannelIdentity.value?.id ?? '',
			channelName: settings.value?.overlaysKappagen?.channel?.login ?? '',
			emotes: {
				ffz: settings.value?.overlaysKappagen?.emotes.ffzEnabled,
				bttv: settings.value?.overlaysKappagen?.emotes.bttvEnabled,
				sevenTv: settings.value?.overlaysKappagen?.emotes.sevenTvEnabled,
			},
			onMessage: onTwitchMessage,
		}
	})

	const { destroy: destroyTwitchChat } = useChatTmi(twitchChatSettings)

	watch([kickChannelIdentity, () => settings.value?.overlaysKappagen.emotes.sevenTvEnabled], ([identity, sevenTvEnabled]) => {
		if (!identity?.id || !sevenTvEnabled) return
		fetchSevenTvEmotes(identity.id, 'kick')
	})

	function randomAnimation(): KappagenAnimations | undefined {
		if (!settings.value?.overlaysKappagen) return

		const enabledAnimations = settings.value?.overlaysKappagen.animations.filter(
			(animation) => animation.enabled
		)

		const index = Math.floor(Math.random() * enabledAnimations.length)
		const randomed = enabledAnimations[index]
		if (!randomed) return { style: 'Confetti', prefs: undefined, count: 150 } as KappagenAnimations

		const splittedAnimationStyle = randomed.style.toLowerCase().split('_')

		const normalizedStyleName = splittedAnimationStyle
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join('')

		return {
			style: normalizedStyleName as KappagenAnimations['style'],
			prefs: randomed.prefs
				? {
						message: randomed.prefs.message ?? [],
						time: randomed.prefs.time ?? 5,
						size: randomed.prefs.size ?? 1,
						speed: randomed.prefs.speed ?? 1,
						faces: randomed.prefs.faces ?? false,
						center: randomed.prefs.center ?? false,
						avoidMiddle: false,
					}
				: undefined,
			count: randomed.count ?? 150,
		} as KappagenAnimations
	}

	watch(eventsData, (event) => {
		if (!event?.twirEvents.baseInfo || !settings.value?.overlaysKappagen) return

		if (
			!settings.value.overlaysKappagen.events.some(
				(e) => e.event === event.twirEvents.baseInfo.type && e.enabled
			)
		) {
			return
		}

		const generatedEmotes = emotesBuilder.buildKappagenEmotes([])

		const animation = randomAnimation()
		if (!animation) return

		unref(instance).playAnimation(generatedEmotes, animation)
	})

	watch(triggerKappagenData, (v) => {
		if (!v?.overlaysKappagenTrigger) return

		const data = v.overlaysKappagenTrigger

		const emotesList: Record<string, string[]> = {}
		const emoteUrls: Record<string, string> = {}
		if (data.emotes) {
			for (const emote of data.emotes) {
				const emoteKey = emote.url ? `${emote.id}\u0000${emote.url}` : emote.id
				emotesList[emoteKey] = [...(emotesList[emoteKey] ?? []), ...emote.positions]
				if (emote.url) {
					emoteUrls[emoteKey] = emote.url
				}
			}
		}

		const chunks = makeMessageChunks(data.text, {
			isSmaller: false,
			emotesList,
			emoteUrls,
		})

		const emotesForKappagen = emotesBuilder.buildKappagenEmotes(chunks)
		const animation = randomAnimation()
		if (!animation) return
		unref(instance).playAnimation(emotesForKappagen, animation)
	})

	watch(chatMessagesData, (v) => {
		const chatMessage = v?.overlaysKappagenChatMessages
		if (!chatMessage || chatMessage.platform === 'twitch') return

		const chunks = makeMessageChunks(chatMessage.text, {
			isSmaller: false,
			emotesList: {},
		})

		handleSpawnChunks(chunks)
	})

	function destroy() {
		pauseEvents()
		pauseSettings()
		pauseChatMessages()
		destroyTwitchChat()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectEvents()
		connectSettings()
		connectChatMessages()
	}

	return {
		connect,
		destroy,
		settings,
	}
}
