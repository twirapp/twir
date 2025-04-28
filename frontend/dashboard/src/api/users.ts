import { useQuery } from '@urql/vue'
import { unref } from 'vue'

import type { MaybeRef } from 'vue'

import { graphql } from '@/gql'

export function useChannelUserInfo(userId: MaybeRef<string | null | undefined>) {
	return useQuery({
		query: graphql(`
			query ChannelUserInfo($userId: String!) {
				channelUserInfo(userId: $userId) {
					userId
					twitchProfile {
						id
						login
						displayName
						profileImageUrl
					}
					watchedMs
					messages
					usedEmotes
					usedChannelPoints
					isMod
					isSubscriber
					isVip
					followerSince
				}
			}
	`),
		get variables() {
			return {
				userId: unref(userId)!,
			}
		},
		pause: !!unref(userId),
	})
}
