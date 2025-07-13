import { useSubscription } from '@urql/vue'
import { type MaybeRef, computed, ref, unref, watch } from 'vue'

import type { Buidler } from './use-kappagen-builder.js'
import type { KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { useMessageHelpers } from '@/composables/tmi/use-message-helpers.js'
import { graphql } from '@/gql'

export function useKappagenOverlaySocket(instance: MaybeRef<KappagenMethods>, emotesBuilder: Buidler) {
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

	const { makeMessageChunks } = useMessageHelpers()

	function randomAnimation(): KappagenAnimations | undefined {
		if (!settings.value?.overlaysKappagen) return

		const enabledAnimations = settings.value?.overlaysKappagen.animations.filter(
			(animation) => animation.enabled,
		)

		const index = Math.floor(Math.random() * enabledAnimations.length)
		const randomed = enabledAnimations[index]

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
				(e) => e.event === event.twirEvents.baseInfo.type && e.enabled,
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
		unref(instance).playAnimation(emotesForKappagen, animation)
	})

	function destroy() {
		pauseEvents()
		pauseSettings()
	}

	async function connect(key: string) {
		apiKey.value = key
		connectEvents()
		connectSettings()
	}

	return {
		connect,
		destroy,
		settings,
	}
}
