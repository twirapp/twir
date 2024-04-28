import { useMutation } from '@urql/vue'
import { defineStore } from 'pinia'

import { graphql } from '@/gql/gql.js'

export const useAdminActions = defineStore('api/admin/actions', () => {
	const useMutationDropUserAuthSession = () => useMutation(graphql(`
		mutation DropUserAuthSession($userId: String!) {
			dropUserAuthSession(userId: $userId)
		}
	`))

	const useMutationDropAllAuthSessions = () => useMutation(graphql(`
		mutation DropAllUserAuthSessions {
			dropAllAuthSessions
		}
	`))

	return {
		useMutationDropUserAuthSession,
		useMutationDropAllAuthSessions
	}
})
