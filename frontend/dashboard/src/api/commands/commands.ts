import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllCommandsQuery } from '@/gql/graphql'

import { commandMenuCacheKey } from '@/api/command-menu.js'
import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export const invalidationKey = 'CommandsInvalidateKey'

export type Command = GetAllCommandsQuery['commands'][0]

export const useCommandsApi = createGlobalState(() => {
	const useQueryCommands = () =>
		useQuery({
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
							onlineOnly
							offlineOnly
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
						offlineOnly
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
						roleCooldowns {
							id
							commandId
							roleId
							cooldown
						}
					}
				}
			`),
			context: {
				additionalTypenames: [invalidationKey],
			},
			variables: {},
		})

	const useMutationDeleteCommand = () =>
		useMutation(
			graphql(`
				mutation DeleteCommand($id: UUID!) {
					commandsRemove(id: $id)
				}
			`),
			[invalidationKey, commandMenuCacheKey]
		)

	const useMutationCreateCommand = () =>
		useMutation(
			graphql(`
				mutation CreateCommand($opts: CommandsCreateOpts!) {
					commandsCreate(opts: $opts) {
						id
					}
				}
			`),
			[invalidationKey, commandMenuCacheKey]
		)

	const useMutationUpdateCommand = () =>
		useMutation(
			graphql(`
				mutation UpdateCommand($id: UUID!, $opts: CommandsUpdateOpts!) {
					commandsUpdate(id: $id, opts: $opts)
				}
			`),
			[invalidationKey, commandMenuCacheKey]
		)

	return {
		useQueryCommands,
		useMutationDeleteCommand,
		useMutationCreateCommand,
		useMutationUpdateCommand,
	}
})
