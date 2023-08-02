<script setup lang="ts">
import { IconGripVertical } from '@tabler/icons-vue';
import { useIntervalFn, useLocalStorage } from '@vueuse/core';
import { intervalToDuration } from 'date-fns';
import { NCard, NSkeleton, NSpace } from 'naive-ui';
import { computed, onBeforeUnmount, ref, watchEffect } from 'vue';
import { VueDraggableNext } from 'vue-draggable-next';

import { useDashboardStats } from '@/api/index.js';
import { padTo2Digits } from '@/helpers/convertMillisToTime.js';


const { data: stats, refetch } = useDashboardStats();

const { pause: pauseStatsFetch } = useIntervalFn(refetch, 5000);

const uptime = ref<string | null>(null);
const { pause } = useIntervalFn(async () => {
	if (!stats.value?.startedAt) return;

	const duration = intervalToDuration({
		start: new Date(Number(stats.value?.startedAt)),
		end: new Date(),
	});

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0];
	if (duration.days !== undefined && duration.days != 0) mappedDuration.unshift(duration.days);

	uptime.value = mappedDuration
		.map(v => padTo2Digits(v!))
		.filter(v => typeof v !== 'undefined')
		.join(':');
}, 1000, { immediate: true });

onBeforeUnmount(() => {
	pause();
	pauseStatsFetch();
});

const statsItems = computed(() => {
	const s = stats.value;
	if (!s) return {};
	const u = uptime.value;

	const items: Record<string, any> = {
		Uptime: u ?? 'Offline',
		Viewers: s.viewers,
		Followers: s.followers,
		Messages: s.chatMessages,
		'Used emotes': s.usedEmotes,
		'Requested songs': s.requestedSongs,
	};

	return items;
});

const localStorageOrder = useLocalStorage<string[]>('twirDashboardStatsOrder', []);

watchEffect(() => {
	const keys = Object.keys(statsItems.value);
	if (!keys.length) return;
	for (const item of keys) {
		if (!localStorageOrder.value.includes(item)) localStorageOrder.value.push(item);
	}
});
</script>

<template>
	<div style="margin:5px">
		<div style="display: flex; flex-direction: row; flex-wrap: wrap; gap: 5px;">
			<n-card style="flex: 1 1 200px;" :bordered="false" embedded content-style="padding: 5px;">
				<span style="font-size:15px">
					{{ stats?.title || 'Offline' }}
				</span>
			</n-card>

			<n-card style="flex: 1 1 200px;" :bordered="false" embedded content-style="padding: 5px;">
				<span style="font-size:15px">
					{{ stats?.categoryName || 'Offline' }}
				</span>
			</n-card>
		</div>

		<Transition mode="out-in" appear>
			<div
				v-if="!Object.keys(statsItems).length"
				style="display: flex; flex-wrap: wrap; width: 100%; gap: 5px; margin-top: 5px"
			>
				<n-skeleton
					v-for="i of 6"
					:key="i"
					:sharp="false"
					height="74px"
					style="flex:1 1 100px; margin-top: 5px;"
				/>
			</div>


			<vue-draggable-next
				v-else
				v-model="localStorageOrder"
				style="display: flex; flex-wrap: wrap; width: 100%; gap: 5px; margin-top: 5px"
			>
				<transition-group>
					<template
						v-for="(item) of localStorageOrder"
						:key="item"
					>
						<n-card
							v-if="typeof statsItems[item] !== 'undefined'"
							style="flex:1 1 100px; cursor: pointer"
							size="small"
							:bordered="false"
							content-style="padding: 5px;"
							embedded
						>
							<n-space vertical>
								<span style="display: flex;">
									<IconGripVertical style="width: 18px;" /> {{ item }}
								</span>
								<span style="font-size:20px">
									{{ statsItems[item] }}
								</span>
							</n-space>
						</n-card>
					</template>
				</transition-group>
			</vue-draggable-next>
		</Transition>
	</div>
</template>

<!-- <style scoped>
.v-enter-active,
.v-leave-active {
	transition: all 1s cubic-bezier(0, 0, 0.2, 1);
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
	transform: scale(0.98);
}
</style> -->

<style scoped>
.v-enter-active,
.v-leave-active {
	transition: opacity 0.1s ease-in-out;
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
}</style>
