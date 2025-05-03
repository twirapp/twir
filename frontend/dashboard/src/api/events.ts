import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'

import type { Event as GqlEvent, EventOperation as GqlEventOperation } from '@/gql/graphql'

import { useMutation } from '@/composables/use-mutation'
import { graphql } from '@/gql'

export type Event = GqlEvent
export type EventOperation = GqlEventOperation
export { EventType } from '@/gql/graphql'

const invalidationKey = 'Events'

export const useEventsApi = createGlobalState(() => {
	// Query all events
	const useQueryEvents = () => useQuery({
		query: graphql(`
      query GetEvents {
        events {
          id
          channelId
          type
          rewardId
          commandId
          keywordId
          description
          enabled
          onlineOnly
          operations {
            id
            type
            input
            delay
            repeat
            useAnnounce
            timeoutTime
            timeoutMessage
            target
            enabled
            filters {
              id
              type
              left
              right
            }
          }
        }
      }
    `),
		context: {
			additionalTypenames: [invalidationKey],
		},
	})

	// Query a single event by ID
	const useQueryEventById = (id: string) => useQuery({
		query: graphql(`
      query GetEventById($id: String!) {
        eventById(id: $id) {
          id
          channelId
          type
          rewardId
          commandId
          keywordId
          description
          enabled
          onlineOnly
          operations {
            id
            type
            input
            delay
            repeat
            useAnnounce
            timeoutTime
            timeoutMessage
            target
            enabled
            filters {
              id
              type
              left
              right
            }
          }
        }
      }
    `),
		variables: { id },
		context: {
			additionalTypenames: [invalidationKey],
		},
		pause: true,
	})

	// Create a new event
	const useMutationCreateEvent = () => useMutation(
		graphql(`
      mutation CreateEvent($input: EventCreateInput!) {
        eventCreate(input: $input) {
          id
        }
      }
    `),
		[invalidationKey],
	)

	// Update an existing event
	const useMutationUpdateEvent = () => useMutation(
		graphql(`
      mutation UpdateEvent($id: String!, $input: EventUpdateInput!) {
        eventUpdate(id: $id, input: $input) {
          id
        }
      }
    `),
		[invalidationKey],
	)

	// Delete an event
	const useMutationDeleteEvent = () => useMutation(
		graphql(`
      mutation DeleteEvent($id: String!) {
        eventDelete(id: $id)
      }
    `),
		[invalidationKey],
	)

	// Enable or disable an event
	const useMutationEnableOrDisableEvent = () => useMutation(
		graphql(`
      mutation EnableOrDisableEvent($id: String!, $enabled: Boolean!) {
        eventEnableOrDisable(id: $id, enabled: $enabled) {
          id
          enabled
        }
      }
    `),
		[invalidationKey],
	)

	return {
		useQueryEvents,
		useQueryEventById,
		useMutationCreateEvent,
		useMutationUpdateEvent,
		useMutationDeleteEvent,
		useMutationEnableOrDisableEvent,
	}
})
