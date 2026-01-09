import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'
import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation.ts'

export const useDonatepayIntegration = createGlobalState(() => {
	const useUpdate = () =>
		useMutation(
			graphql(`
				mutation UpdateDonatepayIntegration($apiKey: String!, $enabled: Boolean!) {
					donatePayIntegration(apiKey: $apiKey, enabled: $enabled) {
						apiKey
						enabled
					}
				}
			`),
			[integrationsPageCacheKey]
		)

	return {
		useUpdate,
	}
})
