import { useSubscription } from '@urql/vue'

import { graphql } from '~/gql'

export function useAllChatMessagesSubscription() {
	return useSubscription({
		query: graphql(`
			subscription AdminChatMessageSubscription {
				adminChatMessages {
					id
					channelId
					channelLogin
					channelName
					userID
					userName
					userDisplayName
					userColor
					text
					createdAt
				}
			}
		`),
		get variables() {
			return {}
		},
	})
}
