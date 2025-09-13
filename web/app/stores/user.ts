import { createRequest } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

export const userProfileWithoutDashboards = createRequest(
	graphql(`
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
			}
		}
	`),
	{}
)

export const UserStoreKey = 'auth-store-user'

export const useAuth = defineStore('auth-store', () => {
	const router = useRouter()

	const {
		data,
		executeQuery: fetchUser,
		fetching,
	} = useQuery({
		query: userProfileWithoutDashboards.query,
		context: {
			key: userProfileWithoutDashboards.key,
			additionalTypenames: [UserStoreKey],
		},
		variables: {},
	})

	const userWithoutDashboards = computed(() => data.value?.authenticatedUser)

	const { executeMutation: executeLogout } = useMutation(
		graphql(`
			mutation userLogout {
				logout
			}
		`)
	)

	async function logout(withRedirect = false) {
		await executeLogout({}, { additionalTypenames: [UserStoreKey] })
		if (withRedirect) {
			await router.push('/')
		}
	}

	const redirectTo = computed(() => {
		const currentRoute = router.currentRoute.value

		const isPublic = currentRoute.matched.at(0)?.path.startsWith('/p/:channelName()')
		const isPaste = currentRoute.matched.at(0)?.path.startsWith('/h')

		if (isPublic || isPaste) {
			return currentRoute.fullPath
		} else {
			return '/dashboard'
		}
	})

	const { executeQuery: fetchAuthLink } = useQuery({
		query: graphql(`
			query AuthLink($redirectTo: String!) {
				authLink(redirectTo: $redirectTo)
			}
		`),
		get variables() {
			return { redirectTo: redirectTo.value }
		},
		pause: true,
	})

	async function getUserDataWithoutDashboards() {
		const { data } = await fetchUser()
		return data.value?.authenticatedUser
	}

	async function login() {
		const { data } = await fetchAuthLink()
		if (!data.value) return

		window.location.replace(data.value.authLink)
	}

	return {
		userWithoutDashboards,
		isLoading: fetching,

		getUserDataWithoutDashboards,
		logout,
		login,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useAuth, import.meta.hot))
}
