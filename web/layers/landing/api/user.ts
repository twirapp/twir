import { createRequest } from '@urql/vue'

import { graphql } from '~/gql'

export const profileQuery = createRequest(graphql(`
	query AuthenticatedUser {
		authenticatedUser {
			id
			isBotAdmin
			isBanned
			isEnabled
			isBotModerator
			hideOnLandingPage
			botId
			apiKey
			twitchProfile {
				description
				displayName
				login
				profileImageUrl
			}
			selectedDashboardId
			selectedDashboardTwitchUser {
				login
				displayName
				profileImageUrl
			}
			availableDashboards {
				id
				flags
				twitchProfile {
					login
					displayName
					profileImageUrl
				}
			}
		}
	}
`), {})

export const userInvalidateQueryKey = 'UserInvalidateQueryKey'

export function useProfile() {
	return useQuery({
		query: profileQuery.query,
		variables: {},
		context: {
			key: profileQuery.key,
			additionalTypenames: [userInvalidateQueryKey],
		},
	})
}

export function useAuthLink(redirectTo: MaybeRef<string>) {
	return useQuery({
		query: graphql(`
			query AuthLink($redirectTo: String!) {
				authLink(redirectTo: $redirectTo)
			}
		`),
		get variables() {
			return {
				redirectTo: unref(redirectTo),
			}
		},
	})
}
