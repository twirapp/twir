import { createGlobalState } from '@vueuse/core'

import { integrationsPageCacheKey, useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

export const useStreamElementsIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()

	const postCode = useMutation(
		graphql(`
			mutation StreamElementsPostCode($input: IntegrationOAuthCodeInput!) {
				streamelementsPostCode(input: $input)
			}
		`),
		[integrationsPageCacheKey]
	)
	const logout = useMutation(
		graphql(`mutation StreamElementsLogout { streamelementsLogout }`),
		[integrationsPageCacheKey]
	)
	const importCommands = useMutation(
		graphql(`
			mutation StreamElementsImportCommands {
				streamelementsImportCommands {
					importedCount failedCount failures { name reason }
				}
			}
		`),
		['commands']
	)
	const importTimers = useMutation(
		graphql(`
			mutation StreamElementsImportTimers {
				streamelementsImportTimers {
					importedCount failedCount failures { name reason }
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
