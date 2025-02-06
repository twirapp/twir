import { createRequest } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

const authedUserQuery = createRequest(graphql(`
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

export const useUser = defineStore('authed-user', () => {
	const { data, executeQuery, fetching, error } = useQuery({
		query: authedUserQuery.query,
		variables: {},
		context: {
			key: authedUserQuery.key,
		},
	})

	const user = computed(() => data.value?.authenticatedUser)

	async function refetch() {
		await executeQuery()
	}

	async function logout() {

	}

	return {
		user,
		isLoading: readonly(fetching),
		error: readonly(error),

		refetch,
		logout,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useUser, import.meta.hot))
}
