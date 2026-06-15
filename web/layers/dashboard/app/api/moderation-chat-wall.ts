import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { ChatWallListQuery } from '~/gql/graphql.js'
import type { MaybeRef } from 'vue'

import { useMutation } from '~~/layers/dashboard/app/composables/use-mutation'
import { graphql } from '~/gql/gql.js'

export type ChatWall = ChatWallListQuery['chatWalls'][0]

export const useModerationChatWall = createGlobalState(() => {
	const useList = () => useQuery({
		query: graphql(`
			query ChatWallList {
				chatWalls {
					id
					phrase
					enabled
					action
					durationSeconds
					timeoutDurationSeconds
					affectedMessages
					createdAt
					updatedAt
				}
			}
		`),
	})

	const useLogs = (logId: MaybeRef<string>) => useQuery({
		query: graphql(`
			query ChatWallLogs($logId: String!) {
				chatWallLogs(id: $logId) {
					id
					userId
					twitchProfile {
						login
						displayName
						profileImageUrl
					}
					text
					createdAt
				}
			}
		`),
		get variables() {
			return {
				logId: unref(logId),
			}
		},
		pause: true,
	})

	const useSettings = () => useQuery({
		query: graphql(`
			query ChatWallSettings {
				chatWallSettings {
					muteSubscribers
					muteVips
				}
			}
		`),
	})

	const useUpdateSettings = () => useMutation(graphql(`
		mutation ChatWallSettingsUpdate($opts: ChatWallSettingsUpdateInput!) {
			chatWallSettingsUpdate(opts: $opts)
		}
	`))

	return {
		useList,
		useLogs,
		useSettings,
		useUpdateSettings,
	}
})
