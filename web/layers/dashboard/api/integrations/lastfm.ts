import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import { graphql } from '@/gql'
import { integrationsPageCacheKey } from '~~/layers/dashboard/api/integrations/integrations-page.js'

export const useLastfmIntegrationApi = createGlobalState(() => {
	const usePostCodeMutation = () =>
		useMutation(
			graphql(`
				mutation LastfmPostCode($code: String!) {
					lastfmPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const useLogoutMutation = () =>
		useMutation(
			graphql(`
				mutation LastfmLogout {
					lastfmLogout
				}
			`),
			[integrationsPageCacheKey]
		)

	return {
		usePostCodeMutation,
		useLogoutMutation,
	}
})
