import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllCommandsQuery } from '@/gql/graphql'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export const invalidationKey = 'CommandsInvalidateKey'

export type Command = GetAllCommandsQuery['commands'][0]

export const useCommandsApi = createGlobalState(() => {
	const useQueryCommands = () => useQuery({
		query: graphql(`
			query GetAllCommands {
				commands {
					id
					name
					description
					aliases
					responses {
						id
						commandId
						text
						twitchCategoriesIds
						twitchCategories {
							id
							name
							boxArtUrl
						}
					}
					cooldown
					cooldownType
					enabled
					visible
					default
					defaultName
					module
					isReply
					keepResponsesOrder
					deniedUsersIds
					allowedUsersIds
					rolesIds
					onlineOnly
					cooldownRolesIds
					enabledCategories
					requiredWatchTime
					requiredMessages
					requiredUsedChannelPoints
					group {
						id
						name
						color
					}
					groupId
					expiresAt
					expiresType
				}
			}
		`),
		context: {
			additionalTypenames: [invalidationKey],
		},
		variables: {},
	})

	const useMutationDeleteCommand = () => useMutation(
		graphql(`
			mutation DeleteCommand($id: ID!) {
				commandsRemove(id: $id)
			}
		`),
		[invalidationKey],
	)

	const useMutationCreateCommand = () => useMutation(
		graphql(`
			mutation CreateCommand($opts: CommandsCreateOpts!) {
				commandsCreate(opts: $opts)
			}
		`),
		[invalidationKey],
	)

	const useMutationUpdateCommand = () => useMutation(
		graphql(`
			mutation UpdateCommand($id: ID!, $opts: CommandsUpdateOpts!) {
				commandsUpdate(id: $id, opts: $opts)
			}
		`),
		[invalidationKey],
	)

	return {
		useQueryCommands,
		useMutationDeleteCommand,
		useMutationCreateCommand,
		useMutationUpdateCommand,
	}
})
