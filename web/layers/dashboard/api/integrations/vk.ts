import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import { graphql } from '~/gql/gql.js'
import { integrationsPageCacheKey } from '~~/layers/dashboard/api/integrations/integrations-page.js'

export const useVKIntegrationApi = createGlobalState(() => {
	const usePostCodeMutation = () =>
		useMutation(
			graphql(`
				mutation VkPostCode($code: String!) {
					vkPostCode(code: $code)
				}
			`),
			[integrationsPageCacheKey]
		)

	const useLogoutMutation = () =>
		useMutation(
			graphql(`
				mutation VkLogout {
					vkLogout
				}
			`),
			[integrationsPageCacheKey]
		)

	return {
		usePostCodeMutation,
		useLogoutMutation,
	}
})
