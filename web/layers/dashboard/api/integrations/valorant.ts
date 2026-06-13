import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import { graphql } from '@/gql'
import { integrationsPageCacheKey } from '~~/layers/dashboard/api/integrations/integrations-page.js'

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
