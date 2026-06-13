import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'
import { integrationsPageCacheKey } from '~~/layers/dashboard/api/integrations/integrations-page.js'

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
