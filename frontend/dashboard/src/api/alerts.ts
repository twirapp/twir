import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { AlertsGetAllQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

const invalidateKey = 'AlertsInvalidateKey'

export type Alert = Omit<AlertsGetAllQuery['channelAlerts'][number], '__typename'>

export const useAlertsApi = createGlobalState(() => {
	const useAlertsQuery = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidateKey] },
		query: graphql(`
			query AlertsGetAll {
				channelAlerts {
					id
					name
					audioId
					audioVolume
					commandIds
					rewardIds
					greetingsIds
					keywordsIds
				}
				keywords {
					id
					text
				}
				commands {
					id
					name
				}
				greetings {
					id
					twitchProfile {
						displayName
						profileImageUrl
					}
				}
				twitchGetChannelRewards {
					partnerOrAffiliate
					rewards {
						id
						title
					}
				}
			}
		`),
	})

	const useAlertsCreateMutation = () => useMutation(graphql(`
		mutation CreateAlert($opts: ChannelAlertCreateInput!) {
			channelAlertsCreate(input: $opts) {
				id
			}
		}
	`), [invalidateKey])

	const useAlertsUpdateMutation = () => useMutation(graphql(`
		mutation UpdateAlert($id: UUID!, $opts: ChannelAlertUpdateInput!) {
			channelAlertsUpdate(id: $id, input: $opts) {
				id
			}
		}
	`), [invalidateKey])

	const useAlertsDeleteMutation = () => useMutation(graphql(`
		mutation DeleteAlert($id: UUID!) {
			channelAlertsDelete(id: $id)
		}
	`), [invalidateKey])

	return {
		useAlertsQuery,
		useAlertsCreateMutation,
		useAlertsUpdateMutation,
		useAlertsDeleteMutation,
	}
})
