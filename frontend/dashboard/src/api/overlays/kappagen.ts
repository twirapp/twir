import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { computed } from 'vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql'

const KappagenOverlayQuery = graphql(`
	query KappagenOverlayQuery {
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

	return {
		kappagen,
		isLoading,
		isUpdating,
		updateKappagen,
		refetch: executeQuery,
	}
})
