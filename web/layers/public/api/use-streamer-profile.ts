import { graphql } from '~/gql'

export function useStreamerProfile() {
	const router = useRouter()

	return useQuery({
		query: graphql(`
			query StreamerTwitchProfile($userName: String!) {
				twitchGetUserByName(name: $userName) {
					id
					profileImageUrl
					login
					description
					displayName
					notFound
				}
			}
		`),
		variables: {
			get userName() {
				return unref(router.currentRoute.value.params.channelName as string ?? '')
			},
		},
	})
}

export function useStreamerPublicSettings(streamerId: MaybeRef<string>) {
	return useQuery({
		query: graphql(`
			query StreamerPublicSettings($streamerId: String!) {
				userPublicSettings(userId: $streamerId) {
					socialLinks {
						title
						href
					}
					description
				}
			}
		`),
		get variables() {
			return {
				streamerId: unref(streamerId),
			}
		},
	})
}
