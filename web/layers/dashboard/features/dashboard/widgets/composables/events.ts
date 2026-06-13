import { createGlobalState, useLocalStorage } from '@vueuse/core'
import { computed } from 'vue'

import { useDashboardEvents } from '@/api/dashboard'
import { DashboardEventType } from '@/gql/graphql'

export const useEvents = createGlobalState(() => {
	const enabledEventsOptions = [
		{
			label: 'Donations',
			value: DashboardEventType.Donation,
		},
		{
			label: 'Follows',
			value: DashboardEventType.Follow,
		},
		{
			label: 'Raids',
			value: DashboardEventType.Raided,
		},
		{
			label: 'Subscriptions',
			value: DashboardEventType.Subscribe,
		},
		{
			label: 'Resubscriptions',
			value: DashboardEventType.Resubscribe,
		},
		{
			label: 'Sub gifts',
			value: DashboardEventType.Subgift,
		},
		{
			label: 'First user messages',
			value: DashboardEventType.FirstUserMessage,
		},
		{
			label: 'Chat clear',
			value: DashboardEventType.ChatClear,
		},
		{
			label: 'Reward activated',
			value: DashboardEventType.RedemptionCreated,
		},
		{
			label: 'Ban/timeout',
			value: DashboardEventType.ChannelBan,
		},
		{
			label: 'Unban request created',
			value: DashboardEventType.ChannelUnbanRequestCreate,
		},
		{
			label: 'Unban request resolved',
			value: DashboardEventType.ChannelUnbanRequestResolve,
		},
	]
	const enabledEvents = useLocalStorage<DashboardEventType[]>(
		'twirEventsWidgetFilterV3',
		enabledEventsOptions.map((e) => e.value)
	)
	const { events, fetching } = useDashboardEvents()

	const filteredEvents = computed(
		() =>
			events.value?.filter((e) => {
				return enabledEvents.value.includes(e.type)
			}) ?? []
	)

	return {
		events: filteredEvents,
		fetching,
		enabledEvents,
		enabledEventsOptions,
	}
})
