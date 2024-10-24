import { useQuery } from '@urql/vue'
import { unref } from 'vue'

import type { AdminAuditLogsInput } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql'

export function useAdminAuditLogs(input: MaybeRef<AdminAuditLogsInput>) {
	return useQuery({
		query: graphql(`
			query AdminAuditLogs($input: AdminAuditLogsInput!) {
				adminAuditLogs(input: $input) {
					logs {
						id
						createdAt
						objectId
						oldValue
						newValue
						operationType
						system
						user {
							id
							profileImageUrl
							displayName
							login
						}
						channel {
							id
							profileImageUrl
							displayName
							login
						}
					}
					total
				}
			}
		`),
		get variables() {
			return {
				input: unref(input),
			}
		},
	})
}
