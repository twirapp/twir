import { useSubscription } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

export function useAllChatMessagesSubscription() {
	return useSubscription({
		query: graphql(`
				subscription AdminChatMessageSubscription {
					adminChatMessages {
						id
						platform
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
