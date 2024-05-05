import { useQuery } from '@urql/vue'

import type { EmotesStatisticQuery, EmotesStatisticsOpts } from '@/gql/graphql'
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
