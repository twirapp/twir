import { useQuery } from '@urql/vue'

import type { AlertsGetAllQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

const invalidateKey = 'AlertsInvalidateKey'

export type Alert = AlertsGetAllQuery['channelAlerts'][0]

export function useAlertsQuery() {
	return useQuery({
		variables: {},
		context: { additionalTypenames: [invalidateKey] },
		query: graphql(`
			query AlertsGetAll {
				channelAlerts {
					id
					channel_id
					name
					audio_id
					audio_volume
					command_ids
					reward_ids
					greetings_ids
					keywords_ids
				}
			}
		`),
	})
}

export function useAlertsCreateMutation() {
	return useMutation(graphql(`
		mutation CreateAlert($opts: ChannelAlertCreateInput!) {
			channelAlertsCreate(input: $opts) {
				id
			}
		}
	`), [invalidateKey])
}

export function useAlertsUpdateMutation() {
	return useMutation(graphql(`
		mutation UpdateAlert($id: ID!, $opts: ChannelAlertUpdateInput!) {
			channelAlertsUpdate(id: $id, input: $opts) {
				id
			}
		}
	`), [invalidateKey])
}

export function useAlertsDeleteMutation() {
	return useMutation(graphql(`
		mutation DeleteAlert($id: ID!) {
			channelAlertsDelete(id: $id)
		}
	`), [invalidateKey])
}
