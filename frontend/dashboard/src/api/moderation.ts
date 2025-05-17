import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { ModerationSettingsCreateOrUpdateInput, ModerationSettingsItem } from '@/gql/graphql.ts'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'

export const useModerationAvailableLanguages = createGlobalState(() => {
	const query = () => useQuery({
		query: graphql(`
			query ModerationAvailableLanguages {
				moderationLanguagesAvailableLanguages {
					languages {
						name
						iso_639_1
					}
				}
			}
		`),
	})

	return {
		query,
	}
})

export type ModerationItem = ModerationSettingsItem
export type ModerationCreateOrUpdateInput = ModerationSettingsCreateOrUpdateInput

export const useChannelModerationSettingsApi = createGlobalState(() => {
	const cacheKey = 'channelModerationSettings'

	const useSettingsQuery = () => useQuery({
		query: graphql(`
			query ChannelModerationSettings {
				moderationSettings {
					id
					type
					name
					enabled
					banTime
					banMessage
					warningMessage
					checkClips
					triggerLength
					maxPercentage
					denyList
					deniedChatLanguages
					denyListRegexpEnabled
					denyListWordBoundaryEnabled
					denyListSensitivityEnabled
					excludedRoles
					maxWarnings
					oneManSpamMessageMemorySeconds
					oneManSpamMinimumStoredMessages
					createdAt
					updatedAt
				}
			}
		`),
		context: {
			additionalTypenames: [cacheKey],
		},
	})

	const useDelete = () => useMutation(graphql(`
		mutation DeleteChannelModerationSetting($id: UUID!) {
			moderationSettingsDelete(id: $id)
		}
	`), [cacheKey])

	const useUpdate = () => useMutation(graphql(`
		mutation UpdateChannelModerationSetting($id: UUID!, $input: ModerationSettingsCreateOrUpdateInput!) {
			moderationSettingsUpdate(id: $id, input: $input) {
				id
			}
		}
	`))

	const useCreate = () => useMutation(graphql(`
		mutation CreateChannelModerationSetting($input: ModerationSettingsCreateOrUpdateInput!) {
			moderationSettingsCreate(input: $input) {
				id
			}
		}
	`), [cacheKey])

	return {
		useQuery: useSettingsQuery,
		useDelete,
		useUpdate,
		useCreate,
	}
})
