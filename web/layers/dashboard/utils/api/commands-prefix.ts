import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql/index.js'

export const useCommandsPrefixApi = createGlobalState(() => {
	const invalidationKey = 'CommandsPrefixInvalidateKey'

	const usePrefix = () => useQuery({
		query: graphql(`
			query CommandsPrefix {
				channelsCommandsPrefix
			}
		`),
		context: {
			additionalTypenames: [invalidationKey],
		},
	})

	const usePrefixUpdate = () => useMutation(
		graphql(`
			mutation CommandsPrefixUpdate($input: CommandsPrefixUpdateInput!) {
				commandsPrefixUpdate(input: $input)
			}
		`),
		[invalidationKey],
	)

	const usePrefixReset = () => useMutation(
		graphql(`
			mutation CommandsPrefixReset {
				commandsPrefixReset
			}
		`),
		[invalidationKey],
	)

	return {
		usePrefix,
		usePrefixUpdate,
		usePrefixReset,
	}
})
