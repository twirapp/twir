import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey, useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

export const useNightbotIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()

	const postCode = useMutation(
		graphql(`
			mutation NightbotPostCode($input: IntegrationOAuthCodeInput!) {
				nightbotPostCode(input: $input)
			}
		`),
		[integrationsPageCacheKey]
	)

	const logout = useMutation(
		graphql(`
			mutation NightbotLogout {
				nightbotLogout
			}
		`),
		[integrationsPageCacheKey]
	)

	const importCommands = useMutation(
		graphql(`
			mutation NightbotImportCommands {
				nightbotImportCommands {
					importedCount
					failedCount
					failures { name reason }
				}
			}
		`),
		['commands']
	)

	const importTimers = useMutation(
		graphql(`
			mutation NightbotImportTimers {
				nightbotImportTimers {
					importedCount
					failedCount
					failures { name reason }
				}
			}
		`),
		['timers']
	)

	return {
		postCode,
		logout,
		importCommands,
		importTimers,
		broadcastRefresh: integrationsPage.broadcastRefresh,
	}
})
