import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllChatAlertsQuery } from '~/gql/graphql.js'

import { useMutation } from '~/composables/use-mutation'
import { graphql } from '~/gql/gql.js'

export type ChatAlerts = GetAllChatAlertsQuery['chatAlerts']

const invalidationKey = 'ChatAlertsInvalidateKey'

export const useChatAlertsApi = createGlobalState(() => {
	const useChatAlertsQuery = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetAllChatAlerts {
				chatAlerts {
					followers {
						enabled
						messages {
							text
						}
						cooldown
					}
					ban {
						enabled
						ignoreTimeoutFrom
						messages {
							text
							count
						}
						cooldown
					}
					raids {
						enabled
						messages {
							text
							count
						}
						cooldown
					}
					donations {
						enabled
						messages {
							text
							count
						}
						cooldown
					}
					subscribers {
						enabled
						messages {
							text
							count
						}
						cooldown
					}
					cheers {
						enabled
						messages {
							text
							count
						}
						cooldown
					}
					redemptions {
						enabled
						messages {
							text
						}
						cooldown
						ignoredRewardsIds
					}
					firstUserMessage {
						enabled
						messages {
							text
						}
						cooldown
					}
					streamOnline {
						enabled
						messages {
							text
						}
						cooldown
					}
					streamOffline {
						enabled
						messages {
							text
						}
						cooldown
					}
					chatCleared {
						enabled
						messages {
							text
						}
						cooldown
					}
					unbanRequestCreate {
						enabled
						messages {
							text
						}
						cooldown
					}
					unbanRequestResolve {
						enabled
						messages {
							text
						}
						cooldown
					}
					messageDelete {
						enabled
						cooldown
						messages {
							text
						}
					}
				}
			}
		`),
	})

	const useMutationUpdateChatAlerts = () => useMutation(graphql(`
		mutation UpdateChatAlerts($input: ChatAlertsInput!) {
			updateChatAlerts(input: $input) {
				__typename
			}
		}
	`), [invalidationKey])

	return {
		useChatAlertsQuery,
		useMutationUpdateChatAlerts,
	}
})
