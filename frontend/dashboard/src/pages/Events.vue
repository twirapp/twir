<script setup lang='ts'>
import { NCard, NGrid, NGridItem, NSkeleton } from 'naive-ui';

import { useEventsManager } from '@/api/index.js';
import Card from '@/components/events/card.vue';

const eventsManager = useEventsManager();
const { data: eventsList, isLoading } = eventsManager.getAll({});
</script>

<template>
	<Transition appear mode="out-in">
		<n-grid v-if="isLoading" responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="index in 6" :key="index" :span="1">
				<n-card content-style="padding: 0px">
					<n-skeleton height="300px" />
				</n-card>
			</n-grid-item>
		</n-grid>

		<n-grid v-else responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="event of eventsList!.events" :key="event.id">
				{{ event }}
				<card style="height: 300px" :event="event" />
			</n-grid-item>
		</n-grid>
	</Transition>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: all 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
