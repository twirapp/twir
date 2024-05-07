import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllKeywordsQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type KeywordResponse = GetAllKeywordsQuery['keywords'][0]
export type Keyword = Omit<KeywordResponse, 'id' | 'response'> & {
	id?: string
	response: string
}

const invalidateKey = 'KeywordsInvalidateKey'

export const useKeywordsApi = createGlobalState(() => {
	const useQueryKeywords = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidateKey] },
		query: graphql(`
			query GetAllKeywords {
				keywords {
					id
					text
					response
					enabled
					cooldown
					isReply
					isRegularExpression
					usageCount
				}
			}
		`),
	})

	const useMutationCreateKeyword = () => useMutation(graphql(`
		mutation CreateKeyword($opts: KeywordCreateInput!) {
			keywordCreate(opts: $opts) {
				id
			}
		}
	`), [invalidateKey])

	const useMutationUpdateKeyword = () => useMutation(graphql(`
		mutation UpdateKeyword($id: String!, $opts: KeywordUpdateInput!) {
			keywordUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidateKey])

	const useMutationRemoveKeyword = () => useMutation(graphql(`
		mutation RemoveKeyword($id: String!) {
			keywordRemove(id: $id)
		}
	`), [invalidateKey])

	return {
		useQueryKeywords,
		useMutationCreateKeyword,
		useMutationUpdateKeyword,
		useMutationRemoveKeyword,
	}
})
