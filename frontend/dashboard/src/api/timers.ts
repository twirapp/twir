import { useQuery } from '@urql/vue';
import { defineStore } from 'pinia';

import { useMutation } from '@/composables/use-mutation.js';
import { graphql } from '@/gql';
import type { GetAllTimersQuery, TimerCreateInput } from '@/gql/graphql';

export type TimerResponse = GetAllTimersQuery['timers'][0];
export type EditableTimer = TimerCreateInput & { id?: string }

const invalidationKey = 'TimersInvalidateKey';

export const useTimersApi = defineStore('api/timers', () => {
	const useQueryTimers = () => useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetAllTimers {
				timers {
					id
					name
					enabled
					timeInterval
					messageInterval
					responses {
						text
						isAnnounce
					}
				}
			}
		`),
	});

	const useMutationCreateTimer = () => useMutation(graphql(`
		mutation CreateTimer($opts: TimerCreateInput!) {
			timersCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey]);

	const useMutationUpdateTimer = () => useMutation(graphql(`
		mutation UpdateTimer($id: String!, $opts: TimerUpdateInput!) {
			timersUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey]);

	const useMutationRemoveTimer = () => useMutation(graphql(`
		mutation RemoveTimer($id: String!) {
			timersRemove(id: $id)
		}
	`), [invalidationKey]);

	return {
		useQueryTimers,
		useMutationCreateTimer,
		useMutationUpdateTimer,
		useMutationRemoveTimer,
	};
});
