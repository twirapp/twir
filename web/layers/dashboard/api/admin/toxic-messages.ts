import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { AdminToxicMessagesInput } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql'

export const useToxicMessagesAdminApi = createGlobalState(() => {
	const invalidationKey = 'AdminToxicMessagesInvalidateKey'

	const useDataQuery = (input: MaybeRef<AdminToxicMessagesInput>) => useQuery({
		query: graphql(`
			query AdminToxicMessages($input: AdminToxicMessagesInput!) {
				adminToxicMessages(input: $input) {
					total
					items {
						id
						channelId
						channelProfile {
							id
							notFound
							description
							login
							displayName
							profileImageUrl
						}
						replyToUserId
						text
						createdAt
					}
				}
			}
		`),
		get variables() {
			return {
				input: unref(input),
			}
		},
		context: {
			additionalTypenames: [invalidationKey],
		},
	})

	return {
		useDataQuery,
	}
})
