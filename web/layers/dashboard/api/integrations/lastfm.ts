import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'
import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation.ts'

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
