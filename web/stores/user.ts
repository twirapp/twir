import { createRequest } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

export const authedUserQuery = createRequest(graphql(`
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

export const useAuth = defineStore('auth-store', () => {
	const { data, executeQuery: fetchUser, fetching, error } = useQuery({
		query: authedUserQuery.query,
		context: {
			key: authedUserQuery.key,
		},
		variables: {},
	})

	const router = useRouter()

	const user = computed(() => data.value?.authenticatedUser)

	const { executeMutation: executeLogout } = useMutation(graphql(`
		mutation userLogout {
			logout
		}
	`))

	async function logout(withRedirect = false) {
		await executeLogout({})
		if (withRedirect) {
			await router.push('/')
		}
		await fetchUser({ requestPolicy: 'network-only' })
	}

	const { data: authLinkData, executeQuery: fetchAuthLink } = useQuery({
		query: graphql(`
			query AuthLink($redirectTo: String!) {
				authLink(redirectTo: $redirectTo)
			}
		`),
		get variables() {
			return { redirectTo: router.currentRoute.value.fullPath }
		},
	})

	const authLink = computed(() => authLinkData.value?.authLink ?? null)

	return {
		user,
		authLink,
		isLoading: fetching,

		refetch: fetchUser,
		fetchAuthLink,
		logout,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useAuth, import.meta.hot))
}
