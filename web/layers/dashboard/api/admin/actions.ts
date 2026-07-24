import { useMutation } from '@urql/vue'

import { graphql } from '~/gql/gql.js'

export const vkVideoBotSetupBroadcastChannelName = 'vk_video_bot_setup' as const

export function useMutationDropAllAuthSessions() {
	return useMutation(graphql(`
		mutation DropAllUserAuthSessions {
			dropAllAuthSessions
		}
	`))
}

export function useMutationEventSubSubscribe() {
	return useMutation(graphql(`
		mutation EventsubSubscribe($opts: EventsubSubscribeInput!) {
			eventsubSubscribe(opts: $opts)
		}
	`))
}

export function useMutationRescheduleTimers() {
	return useMutation(graphql(`
		mutation RescheduleTimers {
			rescheduleTimers
		}
	`))
}

export function useMutationEventSubInitChannels() {
	return useMutation(graphql(`
		mutation EventsubInitChannels {
			eventsubInitChannels
		}
	`))
}

export function useMutationKickBotSetupLink() {
	return useMutation(graphql(`
		mutation KickBotSetupLink {
			kickBotSetupLink
		}
	`))
}

export function useMutationVKVideoBotSetupLink() {
	return useMutation(graphql(`
		mutation VKVideoBotSetupLink {
			vkVideoBotSetupLink
		}
	`))
}

export function useMutationVKVideoBotSetupStatus() {
	return useMutation(graphql(`
		mutation VKVideoBotSetupStatus {
			vkVideoBotSetupStatus
		}
	`))
}

export function useMutationVKVideoBotSetupComplete() {
	return useMutation(graphql(`
		mutation VKVideoBotSetupComplete($code: String!, $state: String!) {
			vkVideoBotSetupComplete(code: $code, state: $state)
		}
	`))
}
