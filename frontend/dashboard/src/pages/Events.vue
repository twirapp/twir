<script setup lang='ts'>
import { useTimeout } from '@vueuse/core';
import { NCard, NGrid, NGridItem, NSkeleton } from 'naive-ui';

import { useEventsManager } from '@/api/index.js';

const eventsManager = useEventsManager();
const { data: eventsList, isLoading } = eventsManager.getAll({});

const t = useTimeout(3000, { immediate: true });

</script>

<template>
	<Transition appear mode="out-in">
		<n-grid v-if="!t" responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="index in 6" :key="index" :span="1">
				<n-card content-style="padding: 0px">
					<n-skeleton height="300px" />
				</n-card>
			</n-grid-item>
		</n-grid>

		<n-grid v-else responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="event of eventsList?.events" :key="event.id">
				<n-card>
					qweq
				</n-card>
			</n-grid-item>
		</n-grid>
	</Transition>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
