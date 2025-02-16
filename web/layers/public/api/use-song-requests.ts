import { graphql } from '~/gql'

export function useSongRequests() {
	const currentChannelId = useCurrentChannelId()

	return useQuery({
		query: graphql(`
			query ChannelPublicSongRequests($channelId: String!) {
				songRequestsPublicQueue(channelId: $channelId) {
					userId
					createdAt
					durationSeconds
					songLink
					title
					twitchProfile {
						id
						login
						displayName
						profileImageUrl
					}
				}
			}
		`),
		get variables() {
			return {
				channelId: unref(currentChannelId) ?? '',
			}
		},
	})
}
