import { useQuery } from '@urql/vue'

import type { EmotesStatisticsOpts } from '@/gql/graphql'
import type { Ref } from 'vue'

import { graphql } from '@/gql'

export function useEmotesStatisticQuery(opts: Ref<EmotesStatisticsOpts>) {
	return useQuery({
		query: graphql(`
		query EmotesStatistic($opts: EmotesStatisticsOpts!) {
			emotesStatistics(opts: $opts) {
				emotes {
					emoteName
					lastUsedAt
					usages
					last24HourUsages {
						count
						usedAt
					}
				}
				total
			}
		}
	`),
		get variables() {
			return {
				opts: opts.value
			}
		}
	})
}
