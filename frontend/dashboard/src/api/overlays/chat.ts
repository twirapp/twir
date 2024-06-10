import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation'
import { graphql } from '@/gql'

const cacheKey = ['chatOverlays']

export const useChatOverlayApi = createGlobalState(() => {
	const useOverlaysQuery = () => useQuery({
		query: graphql(`
			query UseOverlaysData {
				chatOverlays {
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
		context: {
			additionalTypenames: cacheKey,
		},
		variables: {},
	})

	const useOverlayDelete = () => useMutation(
		graphql(`
			mutation DeleteOverlay($id: String!) {
				chatOverlayDelete(id: $id)
			}
		`),
		cacheKey,
	)

	const useOverlayCreate = () => useMutation(
		graphql(`
			mutation CreateOverlay($input: ChatOverlayMutateOpts!) {
				chatOverlayCreate(opts: $input)
			}
		`),
		cacheKey,
	)

	const useOverlayUpdate = () => useMutation(
		graphql(`
			mutation UpdateOverlay($id: String!, $input: ChatOverlayMutateOpts!) {
				chatOverlayUpdate(id: $id, opts: $input)
			}
		`),
		cacheKey,
	)

	return {
		useOverlaysQuery,
		useOverlayDelete,
		useOverlayCreate,
		useOverlayUpdate,
	}
})
