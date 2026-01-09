import { integrationsPageCacheKey } from '#layers/dashboard/api/integrations/integrations-page.ts'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~/composables/use-mutation.ts'
import { graphql } from '~/gql'

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
