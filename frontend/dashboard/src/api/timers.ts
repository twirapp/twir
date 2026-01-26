import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { GetAllTimersQuery, TimerCreateInput } from '@/gql/graphql.js'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql/gql.js'

export type TimerResponse = GetAllTimersQuery['timers'][0]
export type EditableTimer = TimerCreateInput & { id?: string }

const invalidationKey = 'TimersInvalidateKey'

export const useTimersApi = createGlobalState(() => {
	const useQueryTimers = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetAllTimers {
				timers {
					id
					name
					enabled
					offlineEnabled
					timeInterval
					messageInterval
					responses {
						text
						isAnnounce
						count
						announceColor
					}
				}
			}
		`),
	})

	const useMutationCreateTimer = () => useMutation(graphql(`
		mutation CreateTimer($opts: TimerCreateInput!) {
			timersCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationUpdateTimer = () => useMutation(graphql(`
		mutation UpdateTimer($id: UUID!, $opts: TimerUpdateInput!) {
			timersUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationRemoveTimer = () => useMutation(graphql(`
		mutation RemoveTimer($id: UUID!) {
			timersRemove(id: $id)
		}
	`), [invalidationKey])

	return {
		useQueryTimers,
		useMutationCreateTimer,
		useMutationUpdateTimer,
		useMutationRemoveTimer,
	}
})
