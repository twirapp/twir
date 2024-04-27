import { useQuery } from '@urql/vue';
import type { MaybeRef } from 'vue';
import { unref } from 'vue';

import { useMutation } from '@/composables/use-mutation';
import { graphql } from '@/gql';

const userProfileCacheKey = 'userProfile';

export const useUserProfile = () => useQuery({
	query: graphql(`
		query UserProfile {
			authenticatedUser {
				twitchProfile {
					login
					profileImageUrl
					displayName
				}
			}
		}
	`),
	variables: {},
	context: {
		additionalTypenames: [userProfileCacheKey],
	},
});

export const useLogout = () => useMutation(graphql(`
	mutation Logout {
		logout
	}
`), [userProfileCacheKey]);

export const useLoginLink = (redirectTo: MaybeRef<string>) => useQuery({
	query: graphql(`
		query LoginLink($redirectTo: String!) {
			authLink(redirectTo: $redirectTo)
		}
	`),
	variables: {
		get redirectTo() {
			return unref(redirectTo);
		},
	},
});
