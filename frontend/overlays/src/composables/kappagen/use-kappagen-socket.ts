import { useSubscription } from '@urql/vue'
import { ref, watch } from 'vue'

import { useKappagenEmotesBuilder } from './use-kappagen-builder.js'
import type { KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { useMessageHelpers } from '@/composables/tmi/use-message-helpers.js'
import { graphql } from '@/gql'

export const useKappagenOverlaySocket = (instance: MaybeRef<KappagenMethods>) => {
	const apiKey = ref<string>('')

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
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: true,
	})
	const {
		data: settings,
		executeSubscription: connectSettings,
		pause: pauseSettings,
		resume: resumeTrigger,
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
				}
			}
		`),
		get variables() {
			return {
				apiKey: apiKey.value,
			}
		},
		pause: true,
	})
	const { data: triggerKappagenData } = useSubscription({
		query: graphql(`
			subscription KappagenTrigger($apiKey: String!) {
				overlaysKappagenTrigger(apiKey: $apiKey) {
					text
					emotes {
						id
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
	})

	const { makeMessageChunks } = useMessageHelpers()
	const emotesBuilder = useKappagenEmotesBuilder()

	function randomAnimation(): KappagenAnimations | undefined {
		if (!settings.value?.overlaysKappagen) return

		const enabledAnimations = settings.value?.overlaysKappagen.animations.filter(
			(animation) => animation.enabled
		)

		const index = Math.floor(Math.random() * enabledAnimations.length)
		const randomed = enabledAnimations[index]

		const splittedAnimationStyle = randomed.style.toLowerCase().split('_')

		const normalizedStyleName = splittedAnimationStyle
			.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
			.join('')

		return {
			style: normalizedStyleName,
			prefs: randomed.prefs
				? {
						message: randomed.prefs.message,
						time: randomed.prefs.time,
						size: randomed.prefs.size,
						speed: randomed.prefs.speed,
						faces: randomed.prefs.faces,
						center: randomed.prefs.center,
						avoidMiddle: false,
					}
				: null,
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

		const generatedEmotes = emotesBuilder.buildKappagenEmotes([])

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

		const emotesForKappagen = emotesBuilder.buildKappagenEmotes(chunks)
		const animation = randomAnimation()
		if (!animation) return
		options.playAnimation(emotesForKappagen, animation)
	})

	function destroy() {
		pauseEvents()
		pauseSettings()
	}

	async function connect(key: string) {
		apiKey.value = key
		await connectEvents()
		await connectSettings()
	}

	return {
		connect,
		destroy,
		settings,
	}
}
