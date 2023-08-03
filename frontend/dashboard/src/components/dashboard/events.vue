<script lang="ts" setup>
import { EventType } from '@twir/grpc/generated/api/api/dashboard';
import { useIntervalFn, useLocalStorage } from '@vueuse/core';
import { NCard, NScrollbar } from 'naive-ui';
import { computed } from 'vue';

import Donate from './events/donate.vue';
import FirstUserMessage from './events/firstUserMessage.vue';
import Follow from './events/follow.vue';
import Raid from './events/raid.vue';
import ReSubscribe from './events/resubscribe.vue';
import SubGift from './events/subgift.vue';
import Subscribe from './events/subscribe.vue';
import { usePositions } from './positions.js';

import { useDashboardEvents } from '@/api/index.js';

const { data: events, isLoading, refetch } = useDashboardEvents();

useIntervalFn(refetch, 1000);
const enabledEvents = useLocalStorage<number[]>('twirEventsWidgetFilter', Object.values(EventType).filter(t => typeof t === 'number') as number[]);
const filteredEvents = computed(() => events.value?.events.filter(e => {
	console.log(e.type, enabledEvents.value, Object.values(EventType));
	return enabledEvents.value.includes(e.type);
}) ?? []);

const positions = usePositions();
const eventsHeight = computed(() => positions.value.events.height);
</script>

<template>
	<n-card
		title="Events"
		:content-style="{ padding: isLoading ? '10px' : '0px', height: `${eventsHeight-50}px` }"
		:segmented="{
			content: true,
			footer: 'soft'
		}"
		header-style="padding: 5px;"
		:style="{ width: '100%', height: '100%' }"
	>
		<n-scrollbar trigger="none" :style="{ 'max-height': `${eventsHeight - 25}px` }">
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
				</template>
			</TransitionGroup>
		</n-scrollbar>
	</n-card>
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
