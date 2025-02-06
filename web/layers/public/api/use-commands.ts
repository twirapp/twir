import { graphql } from '~/gql'

export function useCommands() {
	const currentStreamerId = useCurrentChannelId()

	return useQuery({
		query: graphql(`
			query PublicCommands($channelId: ID!) {
				commandsPublic(channelId: $channelId) {
					name
					description
					cooldown
					aliases
					cooldownType
					module
					permissions {
						name
						type
					}
					responses
					groupId
					group {
						id
						name
						color
					}
				}
			}
		`),
		get variables() {
			return {
				channelId: currentStreamerId.value ?? '',
			}
		},
	})
}
