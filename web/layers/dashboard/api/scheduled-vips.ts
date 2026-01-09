import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllScheduledVipsQuery } from '@/gql/graphql.ts'

import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation'

const invalidateKey = 'ScheduledVipsInvalidateKey'

export type ScheduledVip = GetAllScheduledVipsQuery['scheduledVips'][0]

export const useScheduledVipsApi = createGlobalState(() => {
	const useQueryScheduledVips = () =>
		useQuery({
			variables: {},
			context: { additionalTypenames: [invalidateKey] },
			query: graphql(`
			query GetAllScheduledVips {
				scheduledVips {
					id
					userID
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
					channelID
					createdAt
					removeAt
				}
			}
		`),
		})

	const useMutationCreateScheduledVip = () =>
		useMutation(
			graphql(`
			mutation CreateScheduledVip($input: ScheduledVipsCreateInput!) {
				scheduledVipsCreate(input: $input)
			}
		`),
			[invalidateKey]
		)

	const useMutationRemoveScheduledVip = () =>
		useMutation(
			graphql(`
			mutation RemoveScheduledVip($id: String!, $input: ScheduledVipsRemoveInput!) {
				scheduledVipsRemove(id: $id, input: $input)
			}
		`),
			[invalidateKey]
		)

	const useMutationUpdateScheduledVip = () =>
		useMutation(
			graphql(`
			mutation UpdateScheduledVip($id: String!, $input: ScheduledVipsUpdateInput!) {
				scheduledVipsUpdate(id: $id, input: $input)
			}
		`),
			[invalidateKey]
		)

	return {
		useQueryScheduledVips,
		useMutationCreateScheduledVip,
		useMutationRemoveScheduledVip,
		useMutationUpdateScheduledVip,
	}
})
