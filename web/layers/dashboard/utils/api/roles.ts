import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql'

const rolesInvalidateKey = 'rolesInvalidateKey'

export const useRoles = createGlobalState(() => {
	const useRolesQuery = () => useQuery({
		query: graphql(`
			query ChannelRoles {
				roles {
				 	id
					type
					name
					permissions
					settings {
						requiredMessages
						requiredUserChannelPoints
						requiredWatchTime
					}
					users {
						id
						displayName
						login
						profileImageUrl
					}
				}
			}
		`),
		variables: {},
		context: {
			additionalTypenames: [rolesInvalidateKey],
		},
	})

	const useRolesDeleteMutation = () => useMutation(
		graphql(`
			mutation DeleteRole($id: UUID!) {
				rolesRemove(id: $id)
			}
		`),
		[rolesInvalidateKey],
	)

	const useRolesCreateMutation = () => useMutation(
		graphql(`
			mutation RolesCreate($opts: RolesCreateOrUpdateOpts!) {
				rolesCreate(opts: $opts)
			}
		`),
		[rolesInvalidateKey],
	)

	const useRolesUpdateMutation = () => useMutation(
		graphql(`
			mutation RolesUpdate($id: UUID!, $opts: RolesCreateOrUpdateOpts!) {
				rolesUpdate(id: $id, opts: $opts)
			}
		`),
		[rolesInvalidateKey],
	)

	return {
		useRolesQuery,
		useRolesDeleteMutation,
		useRolesCreateMutation,
		useRolesUpdateMutation,
	}
})
