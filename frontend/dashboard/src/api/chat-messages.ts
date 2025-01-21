import { useQuery } from '@urql/vue'
import { unref } from 'vue'

import type { ChatMessageInput } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql/gql'

export function useChatMessages(input: MaybeRef<ChatMessageInput>) {
	return useQuery({
		query: graphql(`
			query ChatMessage($input: ChatMessageInput!) {
				chatMessages(input: $input) {
					id
					channelId
					userID
					userName
					userDisplayName
					userColor
					text
					createdAt
					updatedAt
				}
			}
		`),
		get variables() {
			return {
				input: unref(input),
			}
		},
	})
}
