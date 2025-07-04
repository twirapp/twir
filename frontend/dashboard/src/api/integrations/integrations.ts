import { createGlobalState } from '@vueuse/core'
import { useQuery } from '@urql/vue'

import { graphql } from '@/gql'
import { useMutation } from '@/composables/use-mutation.ts'

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
					integrationsDonateStream {
						integrationId
					}
				}
			`),
			context: {
				additionalTypenames: [integrationsCacheKey],
			},
			variables: {},
		})

	const donateStreamPostCode = () =>
		useMutation(
			graphql(`
				mutation DonateStreamPostCode($secret: String!) {
					integrationsDonateStreamPostSecret(input: { secret: $secret })
				}
			`)
		)

	return {
		useQuery: useData,
		donateStreamPostCode,
	}
})
