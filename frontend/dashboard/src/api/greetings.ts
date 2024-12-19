import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllGreetingsQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type Greetings = Omit<GetAllGreetingsQuery['greetings'][0], '__typename'>

const invalidationKey = 'GreetingsInvalidateKey'

export const useGreetingsApi = createGlobalState(() => {
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
		`),
	})

	const useMutationCreateGreetings = () => useMutation(graphql(`
		mutation CreateGreetings($opts: GreetingsCreateInput!) {
			greetingsCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationUpdateGreetings = () => useMutation(graphql(`
		mutation UpdateGreetings($id: UUID!, $opts: GreetingsUpdateInput!) {
			greetingsUpdate(id: UUID, opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationRemoveGreetings = () => useMutation(graphql(`
		mutation RemoveGreetings($id: UUID!) {
			greetingsRemove(id: $id)
		}
	`), [invalidationKey])

	return {
		useQueryGreetings,
		useMutationCreateGreetings,
		useMutationUpdateGreetings,
		useMutationRemoveGreetings,
	}
})
