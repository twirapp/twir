import { createGlobalState } from '@vueuse/core'
import { useMutation } from '@/composables/use-mutation.js'

import { graphql } from '@/gql'
import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.js'

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
