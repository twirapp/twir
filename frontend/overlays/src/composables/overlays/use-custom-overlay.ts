import { useQuery } from '@urql/vue'
import { type MaybeRef, computed, unref } from 'vue'

import { graphql } from '@/gql'

const channelOverlayByIdQuery = graphql(`
	query ChannelOverlayById($id: UUID!) {
		channelOverlayById(id: $id) {
			id
			channelId
			name
			createdAt
			updatedAt
			width
			height
			layers {
				id
				type
				settings {
					htmlOverlayHtml
					htmlOverlayCss
					htmlOverlayJs
					htmlOverlayDataPollSecondsInterval
					imageUrl
				}
				overlayId
				posX
				posY
				width
				height
				rotation
				createdAt
				updatedAt
				periodicallyRefetchData
			}
		}
	}
`)

export const useCustomOverlayById = (id: MaybeRef<string>) => {
	return useQuery({
		query: channelOverlayByIdQuery,
		variables: computed(() => ({
			id: unref(id),
		})),
		pause: computed(() => !unref(id)),
	})
}
