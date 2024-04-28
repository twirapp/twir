import { useMutation } from '@urql/vue';

import { graphql } from '@/gql';

export const useMutationDropUserAuthSession = () => useMutation(graphql(`
	mutation DropUserAuthSession($userId: String!) {
		dropUserAuthSession(userId: $userId)
	}
`));

export const useMutationDropAllAuthSessions = () => useMutation(graphql(`
	mutation DropAllUserAuthSessions {
		dropAllAuthSessions
	}
`));
