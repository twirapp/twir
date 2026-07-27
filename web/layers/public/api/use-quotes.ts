import { graphql } from '~/gql'

export function useQuotes() {
	const currentChannelId = useCurrentChannelId()

	return useQuery({
		query: graphql(`
			query ChannelPublicQuotes($channelId: String!) {
				quotesPublic(channelId: $channelId) {
					number
					text
					creatorName
					gameName
					createdAt
				}
			}
		`),
		get variables() {
			return {
				channelId: unref(currentChannelId) ?? '',
			}
		},
		pause: computed(() => !currentChannelId.value),
	})
}
