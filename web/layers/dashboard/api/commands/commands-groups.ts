import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { invalidationKey as commandsInvalidationKey } from './commands.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

const invalidationKey = 'CommandsGroupsInvalidateKey'

export const useCommandsGroupsApi = createGlobalState(() => {
	const useQueryGroups = () => useQuery({
		query: graphql(`
			query GetAllCommandsGroups {
				commandsGroups {
					id
					name
					color
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [invalidationKey],
		},
	})

	const useMutationDeleteGroup = () => useMutation(
		graphql(`
			mutation DeleteCommandGroup($id: ID!) {
				commandsGroupsRemove(id: $id)
			}
		`),
		[invalidationKey, commandsInvalidationKey],
	)

	const useMutationCreateGroup = () => useMutation(
		graphql(`
			mutation CreateCommandGroup($opts: CommandsGroupsCreateOpts!) {
				commandsGroupsCreate(opts: $opts)
			}
		`),
		[invalidationKey],
	)

	const useMutationUpdateGroup = () => useMutation(
		graphql(`
			mutation UpdateCommandGroup($id: ID!, $opts: CommandsGroupsUpdateOpts!) {
				commandsGroupsUpdate(id: $id,opts: $opts)
			}
		`),
		[invalidationKey, commandsInvalidationKey],
	)

	return {
		useQueryGroups,
		useMutationDeleteGroup,
		useMutationCreateGroup,
		useMutationUpdateGroup,
	}
})
