import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'
import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation.ts'

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
