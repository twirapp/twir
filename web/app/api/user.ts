import { createRequest } from '@urql/vue'
import { computed } from 'vue'

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

export const useProfile = createGlobalState(() => {
	const { data: response, executeQuery, fetching } = useQuery({
		query: profileQuery.query,
		variables: {},
		context: {
			key: profileQuery.key,
			additionalTypenames: [userInvalidateQueryKey],
		},
	})

	const computedUser = computed(() => {
		const user = response.value?.authenticatedUser
		if (!user) return null

		return {
			id: user.id,
			avatar: user.twitchProfile.profileImageUrl,
			login: user.twitchProfile.login,
			displayName: user.twitchProfile.displayName,
			apiKey: user.apiKey,
			isBotAdmin: user.isBotAdmin,
			isEnabled: user.isEnabled,
			isBanned: user.isBanned,
			isBotModerator: user.isBotModerator,
			botId: user.botId,
			selectedDashboardId: user.selectedDashboardId,
			selectedDashboardTwitchUser: user.selectedDashboardTwitchUser,
			hideOnLandingPage: user.hideOnLandingPage,
			availableDashboards: user.availableDashboards,
		}
	})

	return { data: computedUser, executeQuery, isLoading: fetching }
})
