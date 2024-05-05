import { useQuery } from '@urql/vue'

import type {
	EmotesStatisticEmoteDetailedOpts,
	EmotesStatisticQuery,
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
					usagesByUsers {
						date
						user {
							displayName
							profileImageUrl
						}
					}
				}
			}
		`),
	})
}
