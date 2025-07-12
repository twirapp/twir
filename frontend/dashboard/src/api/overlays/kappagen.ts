import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { computed } from 'vue'
import { createGlobalState } from '@vueuse/core'

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
	subscription KappagenOverlaySubscription {
		overlaysKappagen {
			...KappagenOverlaySettings
		}
	}
`)

export const useKappagenApi = createGlobalState(() => {
	const {
		data: kappagenData,
		fetching: isLoading,
		executeQuery,
	} = useQuery({
		query: KappagenOverlayQuery,
	})

	const { executeMutation: updateKappagen, fetching: isUpdating } = useMutation(
		KappagenOverlayUpdateMutation
	)

	const { data: subscriptionData } = useSubscription({
		query: KappagenOverlaySubscription,
	})

	const kappagen = computed(() => {
		return subscriptionData.value?.overlaysKappagen || kappagenData.value?.overlaysKappagen
	})

	const availableAnimations = computed<string[]>(() => {
		return (
			subscriptionData.value?.overlaysKappagenAvailableAnimations ||
			kappagenData.value?.overlaysKappagenAvailableAnimations ||
			[]
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
