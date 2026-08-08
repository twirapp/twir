import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation'
import { graphql } from '~/gql/gql.js'

const invalidationKey = 'DotaSettings'

export const useDotaApi = createGlobalState(() => {
	const useQueryDota = () => useQuery({
		query: graphql(`
			query GetDotaSettings {
				dota {
					id
					channelID
					enabled
					steamAccountId
					steamProfile {
						name
						avatar
						profileUrl
					}
					gsiToken
					gsiConfig
					mmr
					mmrDelta
					sessionWins
					sessionLosses
					predictionSettings {
						enabled
						titleTemplate
						windowSeconds
					}
					chatEvents {
						matchStarted { enabled template cooldown }
						matchEnded { enabled template cooldown }
						roshanKilled { enabled template cooldown }
						aegisPickup { enabled template cooldown }
					}
					commandsSettings {
						mmr
						wl
						lg
						gm
						np
						wp
					}
					createdAt
					updatedAt
				}
			}
		`),
		context: { additionalTypenames: [invalidationKey] },
	})

	const useQueryDotaSteamAuthLink = () => useQuery({
		query: graphql(`
			query GetDotaSteamAuthLink {
				dotaSteamAuthLink
			}
		`),
	})

	const useMutationDotaUpdate = () => useMutation(graphql(`
		mutation UpdateDotaSettings($input: DotaUpdateInput!) {
			dotaUpdate(input: $input) {
				id
				enabled
				mmr
				mmrDelta
				predictionSettings {
					enabled
					titleTemplate
					windowSeconds
				}
				chatEvents {
					matchStarted { enabled template cooldown }
					matchEnded { enabled template cooldown }
					roshanKilled { enabled template cooldown }
					aegisPickup { enabled template cooldown }
				}
				commandsSettings {
					mmr
					wl
					lg
					gm
					np
					wp
				}
			}
		}
	`), [invalidationKey])

	const useMutationDotaSteamLink = () => useMutation(graphql(`
		mutation DotaSteamLink($queryString: String!) {
			dotaSteamLink(queryString: $queryString) {
				id
				steamAccountId
			}
		}
	`), [invalidationKey])

	const useMutationDotaSteamUnlink = () => useMutation(graphql(`
		mutation DotaSteamUnlink {
			dotaSteamUnlink {
				id
				steamAccountId
			}
		}
	`), [invalidationKey])

	const useMutationDotaResetSession = () => useMutation(graphql(`
		mutation DotaResetSession {
			dotaResetSession {
				id
				sessionWins
				sessionLosses
			}
		}
	`), [invalidationKey])

	const useMutationDotaRegenerateGsiToken = () => useMutation(graphql(`
		mutation DotaRegenerateGsiToken {
			dotaRegenerateGsiToken {
				id
				gsiToken
				gsiConfig
			}
		}
	`), [invalidationKey])

	return {
		useQueryDota,
		useQueryDotaSteamAuthLink,
		useMutationDotaUpdate,
		useMutationDotaSteamLink,
		useMutationDotaSteamUnlink,
		useMutationDotaResetSession,
		useMutationDotaRegenerateGsiToken,
	}
})
