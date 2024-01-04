<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue';
import { useIntervalFn } from '@vueuse/core';
import { intervalToDuration } from 'date-fns';
import { useThemeVars, NSkeleton, NText } from 'naive-ui';
import { computed, onBeforeUnmount, ref, h } from 'vue';
import { useI18n } from 'vue-i18n';

import StreamInfoEditor from './components/StreamInfoEditor.vue';

import { useDashboardStats } from '@/api';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';
import { padTo2Digits } from '@/helpers/convertMillisToTime';

const themeVars = useThemeVars();

const { data: stats, refetch, isLoading } = useDashboardStats();

const { pause: pauseStatsFetch } = useIntervalFn(refetch, 5000);

const currentTime = ref(new Date());
const { pause: pauseUptimeInterval } = useIntervalFn(() => {
	currentTime.value = new Date();
}, 1000);

const uptime = computed(() => {
	if (!stats.value?.startedAt) return '00:00:00';

	const duration = intervalToDuration({
		start: new Date(Number(stats.value?.startedAt)),
		end: currentTime.value,
	});

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0];
	if (duration.days !== undefined && duration.days != 0) mappedDuration.unshift(duration.days);

	return mappedDuration
		.map(v => padTo2Digits(v!))
		.filter(v => typeof v !== 'undefined')
		.join(':');
});

onBeforeUnmount(() => {
	pauseStatsFetch();
	pauseUptimeInterval();
});

const { t } = useI18n();

const discrete = useNaiveDiscrete();

function openInfoEditor() {
	discrete.dialog.create({
		showIcon: false,
		content: () => h(StreamInfoEditor, {
			title: stats.value?.title,
			categoryId: stats.value?.categoryId,
			categoryName: stats.value?.categoryName,
		}),
	});
}
</script>

<template>
	<Transition appear mode="out-in">
		<div v-if="isLoading" style="padding-top: 5px; padding-bottom: 5px; width: 100%">
			<n-skeleton width="100%" height="43px" :sharp="false" />
		</div>
		<div v-else class="stats">
			<div class="item stats-uptime" style="cursor: pointer;" @click="openInfoEditor">
				<div
					class="stats-item"
					style="padding-right: 10px;"
				>
					<n-text>
						{{ stats?.title ?? 'No title' }}
					</n-text>
					<n-text>
						{{ stats?.categoryName ?? 'No category' }}
					</n-text>
				</div>
				<IconEdit class="stats-edit-icon" />
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ uptime }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ stats?.viewers ?? 0 }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</n-text>

				<n-text style="font-size: 16px;">
					{{ stats?.followers }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ stats?.chatMessages }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ stats?.subs }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ stats?.usedEmotes }}
				</n-text>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<n-text :depth="3" style="font-size: 11px;">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</n-text>
				<n-text style="font-size: 16px;">
					{{ stats?.requestedSongs }}
				</n-text>
			</div>
		</div>
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

.stats {
	display: flex;
	gap: 12px;
	width: 100%;
	padding-right: 16px;
	padding-left: 16px;
}

.stats-uptime {
	display: flex;
	align-items: center;
}

.stats-edit-icon {
	height: 20px;
	width: 20px;
	cursor: pointer;
}

.item {
	min-width: max-content;
}

.divider {
	border-left: 1px solid v-bind('themeVars.borderColor');
	margin-top: 0.5rem;
	margin-bottom: 0.5rem;
}

.stats-item {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
}
</style>
