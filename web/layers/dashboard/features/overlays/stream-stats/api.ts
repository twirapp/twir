import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import type { OverlaysStreamStatsQuery } from '~/gql/graphql.js'

import { graphql } from '~/gql/gql.js'

export type StreamStatsOverlay = Omit<OverlaysStreamStatsQuery['overlaysStreamStats'], '__typename'>

const invalidationKey = 'StreamStatsOverlayInvalidateKey'

export const useStreamStatsOverlayApi = createGlobalState(() => {
	const useQueryStreamStats = () =>
		useQuery({
			variables: {},
			context: { additionalTypenames: [invalidationKey] },
			query: graphql(`
				query OverlaysStreamStats {
					overlaysStreamStats {
						id
						channelId
						design
						variant
						viewersEnabled
						viewersMode
						platformIconsEnabled
						messagesEnabled
						uptimeEnabled
						subscribersEnabled
						followersEnabled
						viewersColor
						messagesColor
						uptimeColor
						subscribersColor
						followersColor
						counterOrder
						customHtmlEnabled
						customHtml
						customCss
						createdAt
						updatedAt
					}
				}
			`),
		})

	const useMutationUpdateStreamStats = () =>
		useMutation(
			graphql(`
				mutation OverlaysStreamStatsUpdate($input: StreamStatsOverlayUpdateInput!) {
					overlaysStreamStatsUpdate(input: $input) {
						id
						channelId
						design
						variant
						viewersEnabled
						viewersMode
						platformIconsEnabled
						messagesEnabled
						uptimeEnabled
						subscribersEnabled
						followersEnabled
						viewersColor
						messagesColor
						uptimeColor
						subscribersColor
						followersColor
						counterOrder
						customHtmlEnabled
						customHtml
						customCss
						createdAt
						updatedAt
					}
				}
			`),
			[invalidationKey]
		)

	return {
		useQueryStreamStats,
		useMutationUpdateStreamStats,
	}
})
