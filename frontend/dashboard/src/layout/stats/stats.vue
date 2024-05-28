<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { NSkeleton, NText } from 'naive-ui'
import { computed, h, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import StreamInfoEditor from '../stream-info-editor.vue'

import { useRealtimeDashboardStats } from '@/api'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { padTo2Digits } from '@/helpers/convertMillisToTime'

const { stats } = useRealtimeDashboardStats()

const currentTime = ref(new Date())
const { pause: pauseUptimeInterval } = useIntervalFn(() => {
	currentTime.value = new Date()
}, 1000)

const uptime = computed(() => {
	if (!stats.value?.startedAt) return '00:00:00'

	const duration = intervalToDuration({
		start: new Date(stats.value.startedAt),
		end: currentTime.value,
	})

	const mappedDuration = [duration.hours ?? 0, duration.minutes ?? 0, duration.seconds ?? 0]
	if (duration.days !== undefined && duration.days != 0) mappedDuration.unshift(duration.days)

	return mappedDuration
		.map(v => padTo2Digits(v!))
		.filter(v => typeof v !== 'undefined')
		.join(':')
})

onBeforeUnmount(() => {
	pauseUptimeInterval()
})

const { t } = useI18n()

const discrete = useNaiveDiscrete()

function openInfoEditor() {
	discrete.dialog.create({
		showIcon: false,
		content: () => h(StreamInfoEditor, {
			title: stats.value?.title,
			categoryId: stats.value?.categoryId,
			categoryName: stats.value?.categoryName,
		}),
	})
}
</script>

<template>
	<Transition appear mode="out-in">
		<div v-if="!stats" class="py-1 w-full">
			<NSkeleton width="100%" height="43px" :sharp="false" />
		</div>
		<div v-else class="flex gap-3 w-full px-4">
			<div class="item flex items-center cursor-pointer" @click="openInfoEditor">
				<div class="stats-item pr-2.5">
					<NText>
						{{ stats?.title ?? 'No title' }}
					</NText>
					<NText>
						{{ stats?.categoryName ?? 'No category' }}
					</NText>
				</div>
				<IconEdit class="h-5 w-5 cursor-pointer" />
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</NText>
				<NText class="stats-display">
					{{ uptime }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.viewers ?? 0 }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</NText>

				<NText class="stats-display">
					{{ stats?.followers }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.chatMessages }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.subs }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.usedEmotes }}
				</NText>
			</div>

			<div class="divider" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.requestedSongs }}
				</NText>
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

.item {
	@apply min-w-max;
}

.divider {
	@apply my-2 border-l-[color:var(--n-border-color)] border-l border-solid;
}

.stats-item {
	@apply flex flex-col justify-between;
}

.stats-type {
	@apply text-xs;
}

.stats-display {
	@apply text-base tabular-nums;
}
</style>
