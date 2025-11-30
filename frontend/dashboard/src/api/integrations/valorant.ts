import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { useMutation } from '@/composables/use-mutation.ts'

import { graphql } from '@/gql'

const cacheKey = 'valorantIntegrationData'

export const useValorantIntegrationApi = createGlobalState(() => {
	const useData = () =>
		useQuery({
			query: graphql(`
				query ValorantData {
					valorantData {
						enabled
						userName
						avatar
					}
					valorantAuthLink
				}
			`),
			context: { additionalTypenames: [cacheKey] },
		})

	const usePostCodeMutation = () =>
		useMutation(
			graphql(`
				mutation ValorantPostCode($code: String!) {
					valorantPostCode(code: $code)
				}
			`),
			[cacheKey]
		)

	const useLogoutMutation = () =>
		useMutation(
			graphql(`
				mutation ValorantLogout {
					valorantLogout
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
