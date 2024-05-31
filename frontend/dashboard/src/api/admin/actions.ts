import { useMutation } from '@urql/vue'

import { graphql } from '@/gql/gql.js'

export function useMutationDropAllAuthSessions() {
	return useMutation(graphql(`
		mutation DropAllUserAuthSessions {
			dropAllAuthSessions
		}
	`))
}

export function useMutationEventSubSubscribe() {
	return useMutation(graphql(`
		mutation EventsubSubscribe($opts: EventsubSubscribeInput!) {
			eventsubSubscribe(opts: $opts)
		}
	`))
}
