import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { AdminShortUrlsInput, AdminShortUrlsPayload } from '~/gql/graphql'
import type { MaybeRef } from 'vue'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql'

export type AdminShortUrl = AdminShortUrlsPayload['items'][0]

export const useAdminShortUrlsApi = createGlobalState(() => {
	const invalidationKey = 'AdminShortUrlsInvalidateKey'

	const useDataQuery = (input: MaybeRef<AdminShortUrlsInput>) => useQuery({
		query: graphql(`
			query AdminShortUrls($input: AdminShortUrlsInput!) {
				adminShortUrls(input: $input) {
					total
					items {
						id
						link
						userId
						userProfile {
							id
							notFound
							description
							login
							displayName
							profileImageUrl
						}
						views
						createdAt
						updatedAt
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

	const useDeleteMutation = () => useMutation(graphql(`
		mutation DeleteShortUrl($id: String!) {
			adminShortUrlDelete(id: $id)
		}
	`), [invalidationKey])

	const useCreateMutation = () => useMutation(graphql(`
		mutation CreateShortUrl($input: AdminShortUrlCreateInput!) {
			adminShortUrlCreate(input: $input)
		}
	`), [invalidationKey])

	return {
		useQuery: useDataQuery,
		useDeleteMutation,
		useCreateMutation,
	}
})
