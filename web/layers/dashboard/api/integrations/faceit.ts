import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'
import {
	integrationsPageCacheKey,
	useIntegrationsPageData,
} from '~~/layers/dashboard/api/integrations/integrations-page.js'

export const useFaceitIntegration = createGlobalState(() => {
	const faceitBroadcaster = new BroadcastChannel('faceit_channel')
	const integrationsPage = useIntegrationsPageData()

	const postCode = useMutation(
		graphql(`
			mutation FaceitPostCode($code: String!) {
				faceitPostCode(code: $code)
			}
		`),
		[integrationsPageCacheKey]
	)

	const update = useMutation(
		graphql(`
			mutation FaceitUpdate($game: String!) {
				faceitUpdate(game: $game)
			}
		`),
		[integrationsPageCacheKey]
	)

	const logout = useMutation(
		graphql(`
			mutation FaceitLogout {
				faceitLogout
			}
		`),
		[integrationsPageCacheKey]
	)

	faceitBroadcaster.onmessage = (event) => {
		if (event.data !== 'refresh') return
		integrationsPage.refetch()
	}

	function broadcastRefresh() {
		faceitBroadcaster.postMessage('refresh')
	}

	return {
		postCode,
		update,
		logout,
		broadcastRefresh,
	}
})
