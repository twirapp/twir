import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'

export const useDonatepayIntegration = createGlobalState(() => {
	const cacheKey = ['donatepay']

	const data = () =>
		useQuery({
			query: graphql(`
				query DonatepayIntegration {
					donatePayIntegration {
						apiKey
						enabled
					}
				}
			`),
			variables: {},
			context: {
				additionalTypenames: cacheKey,
			},
		})

	const update = () =>
		useMutation(
			graphql(`
				mutation UpdateDonatepayIntegration($apiKey: String!, $enabled: Boolean!) {
					donatePayIntegration(apiKey: $apiKey, enabled: $enabled) {
						apiKey
						enabled
					}
				}
			`),
			cacheKey
		)

	return {
		useQuery: data,
		useUpdate: update,
	}
})
