<script setup lang="ts">
import { IconPencilPlus, IconEyeOff } from '@tabler/icons-vue';
import { useIntervalFn, useLocalStorage } from '@vueuse/core';
import { intervalToDuration } from 'date-fns';
import { GridLayout, GridItem } from 'grid-layout-plus';
import { NButton, NDropdown, useThemeVars, NCard } from 'naive-ui';
import { computed, onBeforeUnmount, ref, watchEffect } from 'vue';
import { useI18n } from 'vue-i18n';

import { useStats } from './stats.js';
import ChannelInfo from './statsChannelInfo.vue';

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
		uptime: u ?? '',
		viewers: s.viewers,
		followers: s.followers,
		messages: s.chatMessages,
		usedEmotes: s.usedEmotes,
		requestedSongs: s.requestedSongs,
		subs: s.subs,
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

const statsWidgets = useStats();
const channelInfoWidget = computed(() => statsWidgets.value.find(v => v.i === 'streamInfo'));
const hideWidget = (key: string | number) => {
	const item = statsWidgets.value.find(i => i.i === key);
	if (!item) return;
	item.visible = false;
};

const dropdownOptions = computed(() => {
	return statsWidgets.value
		.filter((v) => !v.visible)
		.map((v) => ({ label: v.i, key: v.i }));
});

const addWidget = (key: string) => {
	const item = statsWidgets.value.find(v => v.i === key);
	if (!item) return;

	const filteredWidgets = statsWidgets.value.filter(w => w.visible);

	item.visible = true;
  item.y = Math.max(...filteredWidgets.map(w => w.y));
	item.x = Math.max(...filteredWidgets.filter(w => w.y === item.y).map(w => w.x)) + 2;
};

const theme = useThemeVars();

const { t } = useI18n();
</script>

<template>
	<GridLayout
		v-model:layout="statsWidgets"
		:row-height="60"
		:max-rows="3"
		:vertical-compact="false"
	>
		<GridItem
			v-if="channelInfoWidget"
			:key="channelInfoWidget.i"
			:x="channelInfoWidget.x"
			:y="channelInfoWidget.y"
			:w="channelInfoWidget.w"
			:h="channelInfoWidget.h"
			:i="channelInfoWidget.i"
			:min-w="channelInfoWidget.minW ?? 3"
			:min-h="channelInfoWidget.minH ?? 1"
		>
			<ChannelInfo :category-id="stats?.categoryId" :title="stats?.title" :category-name="stats?.categoryName" />
		</GridItem>

		<GridItem
			v-for="widget of statsWidgets.filter(w => w.visible && w.i !== 'streamInfo')"
			:key="widget.i"
			:x="widget.x"
			:y="widget.y"
			:w="widget.w"
			:h="widget.h"
			:i="widget.i"
			:min-w="widget.minW ?? 1"
			:min-h="widget.minH ?? 1"
		>
			<n-card
				v-if="typeof statsItems[widget.i] !== 'undefined'"
				size="small"
				:bordered="false"
				embedded
				content-style="padding: 2px"
				:style="{ 'background-color': theme.actionColor }"
			>
				<n-space vertical>
					<div style="display: flex; justify-content: space-between;">
						<div style="white-space: nowrap; overflow: hidden;text-overflow: ellipsis">
							<span>{{ t(`dashboard.statsWidgets.${widget.i}`) }}</span>
						</div>
						<n-button text @click="hideWidget(widget.i)">
							<IconEyeOff />
						</n-button>
					</div>
					<div style="font-size:20px; white-space: nowrap; overflow: hidden;text-overflow: ellipsis">
						{{ statsItems[widget.i] }}
					</div>
				</n-space>
			</n-card>
		</GridItem>
	</GridLayout>
	<div v-if="dropdownOptions.length" style="padding-right: 10px; padding-top: 10px; padding-bottom: 10px;">
		<n-dropdown placement="left" size="huge" trigger="click" :options="dropdownOptions" @select="addWidget">
			<n-button dashed type="success" style="width: 100%; height: 100%; padding: 5px">
				<IconPencilPlus style="width: 30px; height: 30px" />
			</n-button>
		</n-dropdown>
	</div>
</template>

<style scoped>
.vgl-layout {
  width: 100%;
	user-select: none;
}

.v-enter-active,
.v-leave-active {
	transition: opacity 0.1s ease-in-out;
}

.v-enter-from,
.v-leave-to {
	opacity: 0;
}</style>
