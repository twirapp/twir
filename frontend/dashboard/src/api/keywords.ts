import { useQuery } from '@urql/vue';
import { defineStore } from 'pinia';

import { useMutation } from '@/composables/use-mutation';
import { graphql } from '@/gql';
import type { GetAllKeywordsQuery } from '@/gql/graphql';

export type KeywordResponse = GetAllKeywordsQuery['keywords'][0]
export type Keyword = Required<Omit<KeywordResponse, '__typename'>>

const invalidateKey = 'KeywordsInvalidateKey';

export const useKeywordsApi = defineStore('api/keywords', () => {
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
					usages
				}
			}
		`),
	});

	const useMutationCreateKeyword = () => useMutation(graphql(`
		mutation CreateKeyword($opts: KeywordCreateInput!) {
			keywordCreate(opts: $opts) {
				id
			}
		}
	`), [invalidateKey]);

	const useMutationUpdateKeyword = () => useMutation(graphql(`
		mutation UpdateKeyword($id: String!, $opts: KeywordUpdateInput!) {
			keywordUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidateKey]);

	const useMutationRemoveKeyword = () => useMutation(graphql(`
		mutation RemoveKeyword($id: String!) {
			keywordRemove(id: $id)
		}
	`), [invalidateKey]);

	return {
		useQueryKeywords,
		useMutationCreateKeyword,
		useMutationUpdateKeyword,
		useMutationRemoveKeyword,
	};
});
