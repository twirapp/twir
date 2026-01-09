import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'
import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation.ts'

export const useValorantIntegrationApi = createGlobalState(() => {
	const usePostCodeMutation = () =>
		useMutation(
			graphql(`
				mutation ValorantPostCode($code: String!) {
					valorantPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const useLogoutMutation = () =>
		useMutation(
			graphql(`
				mutation ValorantLogout {
					valorantLogout
				}
			`),
			[integrationsPageCacheKey]
		)

	return {
		usePostCodeMutation,
		useLogoutMutation,
	}
})
