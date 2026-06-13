import { createGlobalState } from '@vueuse/core'
import { useQuery } from '@urql/vue'

import { useMutation } from '@/composables/use-mutation.ts'
import { graphql } from '@/gql'
import { integrationsPageCacheKey } from '@/api/integrations/integrations-page.ts'

export const nightbotCacheKey = 'nightbot'

export const useNightbotIntegration = createGlobalState(() => {
	const nightbotBroadcaster = new BroadcastChannel('nightbot_channel')

	const useAuthLink = (pause = false) =>
		useQuery({
			query: graphql(`
				query NightbotGetAuthLink {
					nightbotGetAuthLink
				}
			`),
			context: {
				additionalTypenames: [nightbotCacheKey],
			},
			variables: {},
			pause,
		})

	const useData = (pause = false) =>
		useQuery({
			query: graphql(`
				query NightbotGetData {
					nightbotGetData {
						userName
						avatar
					}
				}
			`),
			context: {
				additionalTypenames: [nightbotCacheKey],
			},
			variables: {},
			pause,
		})

	const postCode = useMutation(
		graphql(`
			mutation NightbotPostCode($input: NightbotPostCodeInput!) {
				nightbotPostCode(input: $input)
			}
		`),
		[nightbotCacheKey, integrationsPageCacheKey]
	)

	const logout = useMutation(
		graphql(`
			mutation NightbotLogout {
				nightbotLogout
			}
		`),
		[nightbotCacheKey, integrationsPageCacheKey]
	)

	const importCommands = useMutation(
		graphql(`
			mutation NightbotImportCommands {
				nightbotImportCommands {
					importedCount
					failedCount
					failedCommandsNames
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
					failedTimersNames
				}
			}
		`),
		['timers']
	)

	function broadcastRefresh() {
		nightbotBroadcaster.postMessage('refresh')
	}

	return {
		nightbotBroadcaster,
		useAuthLink,
		useData,
		postCode,
		logout,
		importCommands,
		importTimers,
		broadcastRefresh,
	}
})
