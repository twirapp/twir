import { graphql } from '~/gql'

export function useLandingStats() {
	return useQuery({
		query: graphql(`
			query LandingStats {
				twirStats {
					channels
					createdCommands
					messages
					usedCommands
					usedEmotes
					viewers
					streamers {
						id
						followersCount
						isLive
						isPartner
						twitchProfile {
							login
							displayName
							profileImageUrl
						}
					}
				}
			}
		`),
	})
}
