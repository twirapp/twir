import type { Platform, StreamerPublicSettingsQuery } from '~/gql/graphql.js'

import { graphql } from '~/gql'

export const useStreamerProfile = defineStore('streamer-profile', () => {
	const router = useRouter()
	const currentChannelId = useCurrentChannelId()
	const urqlClient = useUrqlClient()

	const { data, executeQuery: executeFetchStreamer } = useQuery({
		query: graphql(`
			query StreamerProfile($userName: String!, $platform: Platform) {
				channelBySlug(channelName: $userName, platform: $platform) {
					id
					hideOnLandingPage
					kickProfile {
						id
						slug
						displayName
						followersCount
						isLive
						profilePicture
					}
					twitchProfile {
						id
						login
						displayName
						profileImageUrl
						description
					}
				}
			}
		`),
		variables: {
			get userName() {
				return unref((router.currentRoute.value.params.channelName as string) ?? '')
			},
			get platform() {
				const platform = String(router.currentRoute.value.params.platform ?? '').toUpperCase()
				return platform === 'TWITCH' || platform === 'KICK' ? (platform as Platform) : undefined
			},
		},
		pause: true,
	})

	const publicProfile = ref<StreamerPublicSettingsQuery>()

	const fetchPublicSettings = (streamerId: string) =>
		urqlClient.query(
			graphql(`
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
			`),
			{ streamerId }
		)

	async function fetchProfile() {
		const { data } = await executeFetchStreamer()
		if (!data.value?.channelBySlug?.id) {
			currentChannelId.value = null
			publicProfile.value = undefined
			return
		}

		const { data: publicData } = await fetchPublicSettings(data.value.channelBySlug.id)
		publicProfile.value = publicData
		currentChannelId.value = data.value.channelBySlug.id ?? null
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
