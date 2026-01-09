import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { ChatMessageInput, ChatMessageQuery } from '~/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '~/gql/gql'

export type ChatMessage = ChatMessageQuery['chatMessages'][number]

export function useChatMessages(input: MaybeRef<ChatMessageInput>) {
	return useQuery({
		query: graphql(`
			query ChatMessage($input: ChatMessageInput!) {
				chatMessages(input: $input) {
					id
					channelId
					channelName
					channelLogin
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
			return {
				input: unref(input),
			}
		},
	})
}

export function useChatMessagesSubscription() {
	return useSubscription({
		query: graphql(`
			subscription ChatMessageSubscription {
				chatMessages {
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

// Global state for chat messages API
export const useChatMessagesApi = createGlobalState(() => {
	const useChatMessagesQuery = (input: MaybeRef<ChatMessageInput>, opts?: { manual?: boolean }) => useQuery({
		query: graphql(`
			query ChatMessage($input: ChatMessageInput!) {
				chatMessages(input: $input) {
					id
					channelId
					channelName
					channelLogin
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
			return {
				input: unref(input),
			}
		},
		pause: opts?.manual ?? false,
	})

	// Subscription for new chat messages
	const subscribeToChatMessages = () => useSubscription({
		query: graphql(`
			subscription SubscribeToChatMessages {
				chatMessages {
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
	})

	return {
		subscribeToChatMessages,
		useQuery: useChatMessagesQuery,
	}
})
