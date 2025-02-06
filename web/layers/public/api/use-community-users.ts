import type { CommunityUsersOpts } from '~/gql/graphql.js'

import { graphql } from '~/gql/gql.js'

export function useCommunityUsers(opts: MaybeRef<CommunityUsersOpts>) {
	return useQuery({
		query: graphql(`
			query GetAllCommunityUsers($opts: CommunityUsersOpts!) {
				communityUsers(opts: $opts) {
					total
					users {
						id
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
						watchedMs
						messages
						usedEmotes
						usedChannelPoints
					}
				}
			}
		`),
		get variables() {
			return {
				opts: unref(opts),
			}
		},
	})
}
