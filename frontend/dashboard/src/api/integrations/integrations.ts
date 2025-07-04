import { createGlobalState } from '@vueuse/core'
import { graphql } from '@/gql'
import { useQuery } from '@urql/vue'

export const integrationsCacheKey = 'integrations'

export const useIntegrations = createGlobalState(() => {
	const useData = () =>
		useQuery({
			query: graphql(`
				query Integrations {
					donatello {
						integrationId
					}
					spotifyData {
						userName
						avatar
					}
					spotifyAuthLink
				}
			`),
			context: {
				additionalTypenames: [integrationsCacheKey],
			},
		})

	return {
		useQuery: useData,
	}
})
