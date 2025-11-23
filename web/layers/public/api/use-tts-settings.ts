import { graphql } from '~/gql'

export function useTtsPublicSettings() {
	const currentChannelId = useCurrentChannelId()

	return useQuery({
		query: graphql(`
			query ChannelPublicTtsSettings($channelId: String!) {
				ttsPublicUsersSettings(channelId: $channelId) {
					userId
					twitchProfile {
						id
						login
						displayName
						profileImageUrl
					}
					pitch
					rate
					voice
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
