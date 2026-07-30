import { graphql } from '~/gql'

export function useLandingStats() {
	return useQuery({
		query: graphql(`
			query LandingStats {
				twirStats {
					channels
				twitchChannels
				kickChannels
				vkChannels
				youtubeChannels
					createdCommands
					messages
					usedCommands
					usedEmotes
					viewers
					shortUrls
					hasteBins
				}
			}
		`),
		variables: {},
	})
}
