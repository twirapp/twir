import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GiveawayFragment } from '@/gql/graphql.ts'

import { graphql } from '@/gql'

graphql(`
	fragment Giveaway on ChannelGiveaway {
		id
		channelId
		createdAt
		updatedAt
		startedAt
		endedAt
		stoppedAt
		keyword
		createdByUserId
		archivedAt
	}
`)

export type Giveaway = GiveawayFragment

export const useGiveawaysApi = createGlobalState(() => {
	const useGiveawaysList = () => useQuery({
		query: graphql(`
			query GiveawaysList {
				giveaways {
					...Giveaway
				}
			}
		`),
	})

	return {
		useGiveawaysList,
	}
})
