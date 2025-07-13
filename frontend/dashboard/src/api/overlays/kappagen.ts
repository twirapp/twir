import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import type { KappagenOverlaySettingsFragment } from '@/gql/graphql.ts'

import { useProfile } from '@/api'
import { graphql } from '@/gql'

graphql(`
	fragment KappagenOverlaySettings on KappagenOverlay {
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
				center
				faces
				size
				speed
				message
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
`)

const KappagenOverlayQuery = graphql(`
	query KappagenOverlayQuery {
		overlaysKappagen {
			...KappagenOverlaySettings
		}
		overlaysKappagenAvailableAnimations
	}
`)

const KappagenOverlayUpdateMutation = graphql(`
	mutation KappagenOverlayUpdate($input: KappagenUpdateInput!) {
		overlaysKappagenUpdate(input: $input) {
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
`)

const KappagenOverlaySubscription = graphql(`
	subscription KappagenOverlaySubscription($apiKey: String!) {
		overlaysKappagen(apiKey: $apiKey) {
			...KappagenOverlaySettings
		}
	}
`)

export const useKappagenApi = createGlobalState(() => {
	const { data: profile } = useProfile()

	const {
		data: kappagenData,
		fetching: isLoading,
		executeQuery,
	} = useQuery({
		query: KappagenOverlayQuery,
	})

	const { executeMutation: updateKappagen, fetching: isUpdating } = useMutation(
		KappagenOverlayUpdateMutation,
	)

	const selectedDashboard = computed(() => {
		return profile.value?.availableDashboards.find(
			(d) => d.id === profile.value?.selectedDashboardId,
		)
	})

	const { data: subscriptionData } = useSubscription({
		query: KappagenOverlaySubscription,
		get variables() {
			return {
				apiKey: profile.value!.apiKey,
			}
		},
		pause: !!selectedDashboard.value,
	})

	const kappagen = computed(() => {
		return subscriptionData.value?.overlaysKappagen || kappagenData.value?.overlaysKappagen
	})

	const availableAnimations = computed<string[]>(() => {
		const data = subscriptionData.value?.overlaysKappagen as KappagenOverlaySettingsFragment | undefined

		return (
			data?.animations?.map((a) => a.style)
			|| kappagenData.value?.overlaysKappagenAvailableAnimations
			|| []
		)
	})

	return {
		kappagen,
		availableAnimations,
		isLoading,
		isUpdating,
		updateKappagen,
		refetch: executeQuery,
	}
})
