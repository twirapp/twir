import { useQuery } from '@urql/vue'
import { defineStore } from 'pinia'

import type { GetAllGreetingsQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type Greetings = GetAllGreetingsQuery['greetings'][0]

const invalidationKey = 'GreetingsInvalidateKey'

export const useGreetingsApi = defineStore('api/greetings', () => {
	const useQueryGreetings = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetAllGreetings {
				greetings {
					id
					userId
					text
					enabled
					isReply
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
				}
			}
		`)
	})

	const useMutationCreateGreetings = () => useMutation(graphql(`
		mutation CreateGreetings($opts: GreetingsCreateInput!) {
			greetingsCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationUpdateGreetings = () => useMutation(graphql(`
		mutation UpdateGreetings($id: String!, $opts: GreetingsUpdateInput!) {
			greetingsUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationRemoveGreetings = () => useMutation(graphql(`
		mutation RemoveGreetings($id: String!) {
			greetingsRemove(id: $id)
		}
	`), [invalidationKey])

	return {
		useQueryGreetings,
		useMutationCreateGreetings,
		useMutationUpdateGreetings,
		useMutationRemoveGreetings
	}
})
