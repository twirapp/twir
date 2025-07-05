import { useSubscription } from '@urql/vue'
import { watch } from 'vue'

import type { Buidler } from './use-kappagen-builder.js'
import type { KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { useMessageHelpers } from '@/composables/tmi/use-message-helpers.js'
import { graphql } from '@/gql'

type Options = Omit<KappagenMethods, 'clear'> & {
	emotesBuilder: Buidler
}

export function useKappagenOverlaySocket(options: Options) {
	const {
		data: eventsData,
		executeSubscription: connectEvents,
		pause: pauseEvents,
	} = useSubscription({
		query: graphql(`
			subscription TwirEvents {
				twirEvents {
					baseInfo {
						channelId
						channelName
						type
					}
				}
			}
		`),
		variables: {},
		pause: true,
	})
	const {
		data: settings,
		executeSubscription: connectSettings,
		pause: pauseSettings,
	} = useSubscription({
		query: graphql(`
			subscription KappagenSettings {
				overlaysKappagen {
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
				}
			}
		`),
		variables: {},
		pause: true,
	})
	const { data: triggerKappagenData } = useSubscription({
		query: graphql(`
			subscription KappagenTrigger {
				overlaysKappagenTrigger {
					text
					emotes {
						id
						positions
					}
				}
			}
		`),
	})

	const { makeMessageChunks } = useMessageHelpers()

	function randomAnimation(): KappagenAnimations | undefined {
		if (!settings.value?.overlaysKappagen) return
		const enabledAnimations = settings.value?.overlaysKappagen.animations.filter(
			(animation) => animation.enabled
		)

		const index = Math.floor(Math.random() * enabledAnimations.length)
		const randomed = enabledAnimations[index]

		return {
			style: randomed.style as KappagenAnimations['style'],
			prefs: {
				message: randomed.prefs.message,
				time: randomed.prefs.time,
				size: randomed.prefs.size,
				speed: randomed.prefs.speed,
				faces: randomed.prefs.faces,
				center: randomed.prefs.center,
				avoidMiddle: false,
			},
			count: randomed.count,
		}
	}

	watch(eventsData, (event) => {
		if (!event.twirEvents.baseInfo || !settings.value?.overlaysKappagen) return

		if (
			!settings.value.overlaysKappagen.events.some(
				(e) => e.event.type === event.twirEvents.baseInfo.type
			)
		) {
			return
		}

		const generatedEmotes = options.emotesBuilder.buildKappagenEmotes([])

		const animation = randomAnimation()
		if (!animation) return

		options.playAnimation(generatedEmotes, animation)
	})

	watch(triggerKappagenData, (v) => {
		if (!v?.overlaysKappagenTrigger) return

		const data = v.overlaysKappagenTrigger

		const emotesList: Record<string, string[]> = {}
		if (data.emotes) {
			for (const emote of data.emotes) {
				emotesList[emote.id] = emote.positions
			}
		}

		const chunks = makeMessageChunks(data.text, {
			isSmaller: false,
			emotesList,
		})

		const emotesForKappagen = options.emotesBuilder.buildKappagenEmotes(chunks)
		const animation = randomAnimation()
		if (!animation) return
		options.playAnimation(emotesForKappagen, animation)
	})

	function destroy() {
		pauseEvents()
		pauseSettings()
	}

	function connect() {
		connectEvents()
		connectSettings()
	}

	return {
		connect,
		destroy,
	}
}
