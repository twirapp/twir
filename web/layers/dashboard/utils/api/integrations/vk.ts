import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~/composables/use-mutation.ts'

import { graphql } from '~/gql'
import { integrationsPageCacheKey } from '#layers/dashboard/api/integrations/integrations-page.ts'

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
