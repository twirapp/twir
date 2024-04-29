import { useQuery } from '@urql/vue'
import { defineStore } from 'pinia'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export const invalidationKey = 'CommandsInvalidateKey'

export const useCommandsApi = defineStore('api/commands', () => {
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
						order
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
				}
			}
		`),
		context: {
			additionalTypenames: [invalidationKey]
		},
		variables: {}
	})

	const useMutationDeleteCommand = () => useMutation(
		graphql(`
			mutation DeleteCommand($id: ID!) {
				commandsRemove(id: $id)
			}
		`),
		[invalidationKey]
	)

	const useMutationCreateCommand = () => useMutation(
		graphql(`
			mutation CreateCommand($opts: CommandsCreateOpts!) {
				commandsCreate(opts: $opts)
			}
		`),
		[invalidationKey]
	)

	const useMutationUpdateCommand = () => useMutation(
		graphql(`
			mutation UpdateCommand($id: ID!, $opts: CommandsUpdateOpts!) {
				commandsUpdate(id: $id, opts: $opts)
			}
		`),
		[invalidationKey]
	)

	return {
		useQueryCommands,
		useMutationDeleteCommand,
		useMutationCreateCommand,
		useMutationUpdateCommand
	}
})
