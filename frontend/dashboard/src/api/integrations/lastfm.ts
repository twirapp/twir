import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { useMutation } from '@/composables/use-mutation.ts'

import { graphql } from '@/gql'

const cacheKey = 'lastfmIntegrationData'

export const useLastfmIntegrationApi = createGlobalState(() => {
	const useData = () =>
		useQuery({
			query: graphql(`
				query LastfmData {
					lastfmData {
						enabled
						userName
						avatar
					}
					lastfmAuthLink
				}
			`),
			context: { additionalTypenames: [cacheKey] },
		})

	const usePostCodeMutation = () =>
		useMutation(
			graphql(`
				mutation LastfmPostCode($code: String!) {
					lastfmPostCode(code: $code)
				}
			`),
			[cacheKey]
		)

	const useLogoutMutation = () =>
		useMutation(
			graphql(`
				mutation LastfmLogout {
					lastfmLogout
				}
			`),
			[cacheKey]
		)

	return {
		useData,
		usePostCodeMutation,
		useLogoutMutation,
	}
})
