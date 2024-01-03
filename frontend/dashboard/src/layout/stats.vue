<script setup lang="ts">
import { IconUser } from '@tabler/icons-vue';
import { useIntervalFn } from '@vueuse/core';
import { intervalToDuration } from 'date-fns';
import { useThemeVars, NSkeleton } from 'naive-ui';
import { computed, onBeforeUnmount, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import StreamInfoEditor from './components/StreamInfoEditor.vue';

import { useDashboardStats } from '@/api';
import { padTo2Digits } from '@/helpers/convertMillisToTime';


const themeVars = useThemeVars();

const { data: stats, refetch, isLoading } = useDashboardStats();

const { pause: pauseStatsFetch } = useIntervalFn(refetch, 5000);

const uptime = computed(() => {
	if (!stats.value?.startedAt) return null;

	const duration = intervalToDuration({
		start: new Date(Number(stats.value?.startedAt)),
		end: new Date(),
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
});

const { t } = useI18n();

const infoEditorOpened = ref(false);
</script>

<template>
	<Transition appear mode="out-in">
		<div v-if="isLoading" class="stats">
			<n-skeleton width="20%" height="25px" :sharp="false" :repeat="6" />
		</div>
		<div v-else class="stats">
			<div class="item" style="position: relative;">
				<div
					class="stats-item"
					style="cursor: pointer;"
					@click="infoEditorOpened = true"
				>
					<span>{{ stats?.title ?? 'No title' }}</span>
					<span>{{ stats?.categoryName ?? 'No category' }}</span>
				</div>
				<div v-if="uptime != null" class="live">
					<div style="display: inline-flex; align-items: center; gap: 2px">
						<span>{{ uptime }}</span>
						|
						<IconUser style="height: 15px; width: 15px;" />
						<span>{{ stats?.viewers }}</span>
					</div>
				</div>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<span class="label">{{ t(`dashboard.statsWidgets.followers`) }}</span>
				<span class="value">{{ stats?.followers }}</span>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<span class="label">{{ t(`dashboard.statsWidgets.messages`) }}</span>
				<span class="value">{{ stats?.chatMessages }}</span>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<span class="label">{{ t(`dashboard.statsWidgets.subs`) }}</span>
				<span class="value">{{ stats?.subs }}</span>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<span class="label">{{ t(`dashboard.statsWidgets.usedEmotes`) }}</span>
				<span class="value">{{ stats?.usedEmotes }}</span>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<span class="label">{{ t(`dashboard.statsWidgets.requestedSongs`) }}</span>
				<span class="value">{{ stats?.requestedSongs }}</span>
			</div>
		</div>
	</Transition>

	<StreamInfoEditor
		:opened="infoEditorOpened"
		:title="stats?.title"
		:categoryId="stats?.categoryId"
		:category-name="stats?.categoryName"
		@close="infoEditorOpened = false"
	/>
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
}

.item {
	padding-right: 10px;
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

.stats-item .label {
	font-size: 0.7rem;
	font-weight: 400;
}

.stats-item .value {
	font-size: 1rem;
}

.live {
	position: absolute;
	right: -15px;
	top: -10px;
	background-color: rgb(220 38 38);
	font-size: .75rem;
	line-height: 1rem;
	font-weight: 600;
	padding-left: 10px;
	padding-right: 10px;
	border-radius: 4px;
}
</style>
