<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { NSkeleton, NText } from 'naive-ui'
import { computed, h, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import StreamInfoEditor from '../stream-info-editor.vue'

import { useRealtimeDashboardStats } from '@/api'
import { Separator } from '@/components/ui/separator'
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
	if (duration.days !== undefined && duration.days !== 0) mappedDuration.unshift(duration.days)

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
		<div v-if="!stats" class="w-full p-2">
			<NSkeleton width="100%" height="45px" :sharp="false" style="border-radius: 10px" />
		</div>

		<div v-else class="flex bg-card gap-3 py-2 px-2 flex-wrap w-full border-b border-b-border justify-center">
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

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</NText>
				<NText class="stats-display">
					{{ uptime }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.viewers ?? 0 }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</NText>

				<NText class="stats-display">
					{{ stats?.followers }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.chatMessages }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.subs }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.usedEmotes }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.requestedSongs }}
				</NText>
			</div>

			<Separator orientation="vertical" class="mx-2" />
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
