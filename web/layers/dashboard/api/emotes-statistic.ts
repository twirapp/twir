import { useQuery } from '@urql/vue'

import type {
	EmotesStatisticEmoteDetailedOpts,
	EmotesStatisticQuery,
	EmotesStatisticsDetailsQuery,
	EmotesStatisticsOpts,
} from '@/gql/graphql'
import type { Ref } from 'vue'

import { graphql } from '@/gql'

export type EmotesStatistics = EmotesStatisticQuery['emotesStatistics']['emotes']

export function useEmotesStatisticQuery(opts: Ref<EmotesStatisticsOpts>) {
	return useQuery({
		get variables() {
			return {
				opts: opts.value,
			}
		},
		query: graphql(`
			query EmotesStatistic($opts: EmotesStatisticsOpts!) {
				emotesStatistics(opts: $opts) {
					emotes {
						emoteName
						totalUsages
						lastUsedTimestamp
						graphicUsages {
							count
							timestamp
						}
					}
					total
				}
			}
		`),
	})
}

export type EmotesStatisticsDetail = EmotesStatisticsDetailsQuery

export function useEmotesStatisticDetailsQuery(opts: Ref<EmotesStatisticEmoteDetailedOpts>) {
	return useQuery({
		get variables() {
			return {
				opts: opts.value,
			}
		},
		query: graphql(`
			query EmotesStatisticsDetails($opts: EmotesStatisticEmoteDetailedOpts!) {
				emotesStatisticEmoteDetailedInformation(opts: $opts) {
					graphicUsages {
						count
						timestamp
					}
					usagesHistory {
						date
						userId
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
					}
					topUsers {
						userId
						count
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
					}
					usagesByUsersTotal
					topUsersTotal
				}
			}
		`),
	})
}
