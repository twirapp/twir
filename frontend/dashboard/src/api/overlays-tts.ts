import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { OverlaysTtsQuery } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type TTSOverlay = Omit<OverlaysTtsQuery['overlaysTTS'], '__typename' | 'channel'>

const invalidationKey = 'TTSOverlayInvalidateKey'

export const useTTSOverlayApi = createGlobalState(() => {
	const useQueryTTS = () =>
		useQuery({
			variables: {},
			context: { additionalTypenames: [invalidationKey] },
			query: graphql(`
				query OverlaysTTS {
					overlaysTTS {
						id
						enabled
						voice
						disallowedVoices
						pitch
						rate
						volume
						doNotReadTwitchEmotes
						doNotReadEmoji
						doNotReadLinks
						allowUsersChooseVoiceInMainCommand
						maxSymbols
						readChatMessages
						readChatMessagesNicknames
						createdAt
						updatedAt
						channelId
					}
				}
			`),
		})

	const useMutationUpdateTTS = () =>
		useMutation(
			graphql(`
				mutation OverlaysTTSUpdate($input: TTSUpdateInput!) {
					overlaysTTSUpdate(input: $input) {
						id
						enabled
						voice
						disallowedVoices
						pitch
						rate
						volume
						doNotReadTwitchEmotes
						doNotReadEmoji
						doNotReadLinks
						allowUsersChooseVoiceInMainCommand
						maxSymbols
						readChatMessages
						readChatMessagesNicknames
						createdAt
						updatedAt
						channelId
					}
				}
			`),
			[invalidationKey]
		)

	const useQueryTTSUsersSettings = () =>
		useQuery({
			variables: {},
			query: graphql(`
				query OverlaysTTSUsersSettings {
					overlaysTTSUsersSettings {
						userId
						twitchProfile {
							id
							login
							displayName
							profileImageUrl
						}
						rate
						pitch
						voice
						isChannelOwner
					}
				}
			`),
		})

	const useMutationDeleteTTSUsersSettings = () =>
		useMutation(
			graphql(`
				mutation OverlaysTTSUsersDelete($userIds: [String!]!) {
					overlaysTTSUsersDelete(userIds: $userIds)
				}
			`),
			[]
		)

	const useQueryTTSGetInfo = () =>
		useQuery({
			variables: {},
			query: graphql(`
				query OverlaysTTSGetInfo {
					overlaysTTSGetInfo {
						voicesInfo {
							key
							info {
								country
								gender
								lang
								name
								no
							}
						}
					}
				}
			`),
		})

	return {
		useQueryTTS,
		useMutationUpdateTTS,
		useQueryTTSUsersSettings,
		useMutationDeleteTTSUsersSettings,
		useQueryTTSGetInfo,
	}
})
