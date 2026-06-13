import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { ChatMessageInput, ChatMessageQuery } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql/gql'

export type ChatMessage = ChatMessageQuery['chatMessages'][number]

export function useChatMessages(input: MaybeRef<ChatMessageInput>) {
	return useQuery({
		// @ts-ignore
		query: graphql(`
			query ChatMessage($input: ChatMessageInput!) {
				chatMessages(input: $input) {
					id
					platform
					channelId
					channelName
					channelLogin
					userID
					userName
					userDisplayName
					userColor
					text
					createdAt
					platform
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
		// @ts-ignore
		query: graphql(`
			subscription ChatMessageSubscription {
				chatMessages {
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
					platform
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
		// @ts-ignore
		query: graphql(`
			query ChatMessage($input: ChatMessageInput!) {
				chatMessages(input: $input) {
					id
					platform
					channelId
					channelName
					channelLogin
					userID
					userName
					userDisplayName
					userColor
					text
					createdAt
					platform
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
		// @ts-ignore
		query: graphql(`
			subscription SubscribeToChatMessages {
				chatMessages {
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
					platform
				}
			}
		`),
	})

	return {
		subscribeToChatMessages,
		useQuery: useChatMessagesQuery,
	}
})
