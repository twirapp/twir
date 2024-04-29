import { useMutation } from '@urql/vue'

import { graphql } from '@/gql/gql.js'

export function useMutationDropAllAuthSessions() {
	return useMutation(graphql(`
		mutation DropAllUserAuthSessions {
			dropAllAuthSessions
		}
	`))
}
