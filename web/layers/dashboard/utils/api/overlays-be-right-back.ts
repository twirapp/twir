import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { OverlaysBeRightBackQuery } from '~/gql/graphql.js'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql/gql.js'

export type BeRightBackOverlay = Omit<
	OverlaysBeRightBackQuery['overlaysBeRightBack'],
	'__typename' | 'channel'
>

const invalidationKey = 'BeRightBackOverlayInvalidateKey'

export const useBeRightBackOverlayApi = createGlobalState(() => {
	const useQueryBeRightBack = () =>
		useQuery({
			variables: {},
			context: { additionalTypenames: [invalidationKey] },
			query: graphql(`
				query OverlaysBeRightBack {
					overlaysBeRightBack {
						id
						text
						late {
							enabled
							text
							displayBrbTime
						}
						backgroundColor
						fontSize
						fontColor
						fontFamily
						createdAt
						updatedAt
						channelId
					}
				}
			`),
		})

	const useMutationUpdateBeRightBack = () =>
		useMutation(
			graphql(`
				mutation OverlaysBeRightBackUpdate($input: BeRightBackUpdateInput!) {
					overlaysBeRightBackUpdate(input: $input) {
						id
						text
						late {
							enabled
							text
							displayBrbTime
						}
						backgroundColor
						fontSize
						fontColor
						fontFamily
						createdAt
						updatedAt
						channelId
					}
				}
			`),
			[invalidationKey]
		)

	return {
		useQueryBeRightBack,
		useMutationUpdateBeRightBack,
	}
})
