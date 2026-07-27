import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllQuotesQuery } from '~/gql/graphql.js'

import { commandMenuCacheKey } from '~~/layers/dashboard/api/command-menu.js'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

export type Quote = Omit<GetAllQuotesQuery['quotes'][0], '__typename'>

const QuotesInvalidateKey = 'QuotesInvalidateKey'

export const useQuotesApi = createGlobalState(() => {
	const useQueryQuotes = () => useQuery({
		variables: {},
		context: { additionalTypenames: [QuotesInvalidateKey] },
		query: graphql(`
			query GetAllQuotes {
				quotes {
					id
					number
					text
					creatorId
					creatorName
					gameId
					gameName
					createdAt
					updatedAt
				}
			}
		`),
	})

	const useMutationCreateQuote = () => useMutation(graphql(`
		mutation QuoteCreate($opts: QuoteCreateInput!) {
			quoteCreate(opts: $opts) {
				id
			}
		}
	`), [QuotesInvalidateKey, commandMenuCacheKey])

	const useMutationUpdateQuote = () => useMutation(graphql(`
		mutation QuoteUpdate($id: UUID!, $opts: QuoteUpdateInput!) {
			quoteUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [QuotesInvalidateKey, commandMenuCacheKey])

	const useMutationRemoveQuote = () => useMutation(graphql(`
		mutation QuoteRemove($id: UUID!) {
			quoteRemove(id: $id)
		}
	`), [QuotesInvalidateKey, commandMenuCacheKey])

	return {
		useQueryQuotes,
		useMutationCreateQuote,
		useMutationUpdateQuote,
		useMutationRemoveQuote,
	}
})
