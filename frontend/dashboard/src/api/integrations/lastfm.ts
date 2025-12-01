import { createGlobalState } from '@vueuse/core'
import { useMutation } from '@/composables/use-mutation.ts'

import { graphql } from '@/gql'
import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'

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
