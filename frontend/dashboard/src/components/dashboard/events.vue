<script lang="ts" setup>
import { IconSettings } from '@tabler/icons-vue'
import { EventType } from '@twir/api/messages/dashboard/dashboard'
import { useIntervalFn, useLocalStorage } from '@vueuse/core'
import { NButton, NPopselect, NScrollbar, NText } from 'naive-ui'
import { computed } from 'vue'

import Card from './card.vue'
import Ban from './events/ban.vue'
import ChatClear from './events/chatClear.vue'
import Donate from './events/donate.vue'
import FirstUserMessage from './events/firstUserMessage.vue'
import Follow from './events/follow.vue'
import Raid from './events/raid.vue'
import RedemptionCreated from './events/redemptionCreated.vue'
import ReSubscribe from './events/resubscribe.vue'
import SubGift from './events/subgift.vue'
import Subscribe from './events/subscribe.vue'

import { useDashboardEvents } from '@/api/index.js'
import UnbanRequestCreated from '@/components/dashboard/events/unban-request-created.vue'
import UnbanRequestResolved from '@/components/dashboard/events/unban-request-resolved.vue'

const { data: events, isLoading, refetch } = useDashboardEvents()
useIntervalFn(refetch, 1000)

const enabledEvents = useLocalStorage<number[]>('twirEventsWidgetFilterv2', Object.values(EventType).filter(t => typeof t === 'number') as number[])
const filteredEvents = computed(() => events.value?.events.filter(e => {
	return enabledEvents.value.includes(e.type)
}) ?? [])

const enabledEventsOptions = [
	{
		label: 'Donations',
		value: 0,
	},
	{
		label: 'Follows',
		value: 1,
	},
	{
		label: 'Raids',
		value: 2,
	},
	{
		label: 'Subscriptions',
		value: 3,
	},
	{
		label: 'Resubscriptions',
		value: 4,
	},
	{
		label: 'Sub gifts',
		value: 5,
	},
	{
		label: 'First user messages',
		value: 6,
	},
	{
		label: 'Chat clear',
		value: 7,
	},
	{
		label: 'Reward activated',
		value: 8,
	},
	{
		label: 'Ban/timeout',
		value: 9,
	},
	{
		label: 'Unban request created',
		value: 10,
	},
	{
		label: 'Unban request resolved',
		value: 11,
	},
]
</script>

<template>
	<Card :content-style="{ padding: isLoading ? '10px' : '0px', height: '80%' }">
		<template #header-extra>
			<NPopselect
				v-model:value="enabledEvents" multiple :options="enabledEventsOptions"
				trigger="click"
			>
				<NButton text>
					<IconSettings />
				</NButton>
			</NPopselect>
		</template>
		<NScrollbar v-if="filteredEvents.length" trigger="none">
			<TransitionGroup name="list">
				<template v-for="(event) of filteredEvents" :key="event.createdAt">
					<Follow
						v-if="event.type === EventType.FOLLOW"
						:created-at="event.createdAt"
						:user-name="event.data!.followUserName"
						:user-display-name="event.data!.followUserDisplayName"
					/>
					<Raid
						v-if="event.type === EventType.RAIDED"
						:created-at="event.createdAt"
						:user-name="event.data!.raidedFromUserName"
						:user-display-name="event.data!.raidedFromDisplayName"
						:viewers="event.data!.raidedViewersCount"
					/>
					<Donate
						v-if="event.type === EventType.DONATION"
						:created-at="event.createdAt"
						:user-name="event.data!.donationUsername"
						:amount="event.data!.donationAmount"
						:message="event.data!.donationMessage"
						:currency="event.data!.donationCurrency"
					/>
					<Subscribe
						v-if="event.type === EventType.SUBSCRIBE"
						:created-at="event.createdAt"
						:user-name="event.data!.subUserName"
						:user-display-name="event.data!.subUserDisplayName"
						:level="event.data!.subLevel"
					/>
					<ReSubscribe
						v-if="event.type === EventType.RESUBSCRIBE"
						:created-at="event.createdAt"
						:user-name="event.data!.reSubUserName"
						:user-display-name="event.data!.reSubUserDisplayName"
						:level="event.data!.reSubLevel"
						:months="event.data!.reSubMonths"
						:streak="event.data!.reSubStreak"
					/>
					<SubGift
						v-if="event.type === EventType.SUBGIFT"
						:created-at="event.createdAt"
						:user-name="event.data!.subGiftTargetUserName"
						:user-display-name="event.data!.subGiftUserDisplayName"
						:level="event.data!.subGiftLevel"
						:target-user-name="event.data!.subGiftTargetUserName"
						:target-user-display-name="event.data!.subGiftTargetUserDisplayName"
					/>
					<FirstUserMessage
						v-if="event.type === EventType.FIRST_USER_MESSAGE"
						:created-at="event.createdAt"
						:user-name="event.data!.firstUserMessageUserName"
						:user-display-name="event.data!.firstUserMessageUserDisplayName"
						:message="event.data!.firstUserMessageMessage"
					/>
					<ChatClear
						v-if="event.type === EventType.CHAT_CLEAR"
						:created-at="event.createdAt"
					/>
					<RedemptionCreated
						v-if="event.type === EventType.REDEMPTION_CREATED"
						:created-at="event.createdAt"
						:title="event.data!.redemptionTitle"
						:input="event.data!.redemptionInput"
						:user-name="event.data!.redemptionUserName"
						:user-display-name="event.data!.redemptionUserDisplayName"
						:cost="event.data!.redemptionCost"
					/>
					<Ban
						v-if="event.type === EventType.CHANNEL_BAN"
						:created-at="event.createdAt"
						:ends-in="event.data!.banEndsInMinutes"
						:moderator-user-login="event.data!.moderatorName"
						:moderator-user-name="event.data!.moderatorDisplayName"
						:reason="event.data!.banReason"
						:user-login="event.data!.bannedUserLogin"
						:user-name="event.data!.bannedUserName"
					/>
					<UnbanRequestCreated
						v-if="event.type === EventType.CHANNEL_UNBAN_REQUEST_CREATE"
						:created-at="event.createdAt"
						:message="event.data!.message"
						:user-login="event.data!.userLogin"
						:user-name="event.data!.userName"
					/>
					<UnbanRequestResolved
						v-if="event.type === EventType.CHANNEL_UNBAN_REQUEST_RESOLVE"
						:created-at="event.createdAt"
						:message="event.data!.message"
						:user-login="event.data!.userLogin"
						:user-name="event.data!.userName"
						:moderator-user-login="event.data!.moderatorName"
						:moderator-user-name="event.data!.moderatorDisplayName"
					/>
				</template>
			</TransitionGroup>
		</NScrollbar>
		<div v-else class="flex items-center justify-center h-full">
			<NText class="text-4xl">
				No events
			</NText>
		</div>
	</Card>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
	transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
	opacity: 0;
	transform: translateY(30px);
}
</style>
