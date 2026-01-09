import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { graphql } from '@/gql/index.js'
import { useMutation } from '~/composables/use-mutation.js'

export const useCommandsPrefixApi = createGlobalState(() => {
	const invalidationKey = 'CommandsPrefixInvalidateKey'

	const usePrefix = () =>
		useQuery({
			query: graphql(`
			query CommandsPrefix {
				channelsCommandsPrefix
			}
		`),
			context: {
				additionalTypenames: [invalidationKey],
			},
		})

	const usePrefixUpdate = () =>
		useMutation(
			graphql(`
			mutation CommandsPrefixUpdate($input: CommandsPrefixUpdateInput!) {
				commandsPrefixUpdate(input: $input)
			}
		`),
			[invalidationKey]
		)

	const usePrefixReset = () =>
		useMutation(
			graphql(`
			mutation CommandsPrefixReset {
				commandsPrefixReset
			}
		`),
			[invalidationKey]
		)

	return {
		usePrefix,
		usePrefixUpdate,
		usePrefixReset,
	}
})
