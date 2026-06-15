import type { StreamerPublicSettingsQuery } from '~/gql/graphql.js'

import { graphql } from '~/gql'

export const useStreamerProfile = defineStore('streamer-profile', () => {
	const router = useRouter()
	const currentChannelId = useCurrentChannelId()
	const urqlClient = useUrqlClient()

	const { data, executeQuery: executeFetchStreamer } = useQuery({
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
		pause: true,
	})

	const publicProfile = ref<StreamerPublicSettingsQuery>()

	const fetchPublicSettings = (streamerId: string) => urqlClient.query(graphql(`
		query StreamerPublicSettings($streamerId: String!) {
			userPublicSettings(userId: $streamerId) {
				channelId
				socialLinks {
					title
					href
				}
				description
			}
		}
	`), { streamerId })

	async function fetchProfile() {
		const { data } = await executeFetchStreamer()
		if (!data.value?.twitchGetUserByName?.id) {
			currentChannelId.value = null
			publicProfile.value = undefined
			return
		}

		const { data: publicData } = await fetchPublicSettings(data.value.twitchGetUserByName.id)
		publicProfile.value = publicData
		currentChannelId.value = publicData?.userPublicSettings?.channelId ?? null
	}

	return {
		profile: data,
		publicProfile,
		fetchProfile,
	}
})

if (import.meta.hot) {
	import.meta.hot.accept(acceptHMRUpdate(useStreamerProfile, import.meta.hot))
}
