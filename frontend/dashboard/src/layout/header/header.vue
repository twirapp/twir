<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { computed, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import StreamInfoEditor from '../stream-info-editor.vue'

import { useRealtimeDashboardStats } from '@/api/dashboard'
import { padTo2Digits } from '@/helpers/convertMillisToTime'
import HeaderProfile from '@/layout/header/header-profile.vue'
import HeaderBotStatus from '@/layout/header/header-bot-status.vue'

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
		.map((v) => padTo2Digits(v!))
		.filter((v) => typeof v !== 'undefined')
		.join(':')
})

onBeforeUnmount(() => {
	pauseUptimeInterval()
})

const { t } = useI18n()

const streamInfoEditorOpen = ref(false)

function openInfoEditor() {
	streamInfoEditorOpen.value = true
}
</script>

<template>
	<div
		class="flex flex-wrap justify-between bg-card w-full px-2 py-1 gap-2 border-b border-b-border"
	>
		<div class="flex flex-wrap md:flex-row flex-col gap-4 py-2">
			<div class="flex items-center cursor-pointer overflow-hidden" @click="openInfoEditor">
				<div class="flex flex-row md:flex-col pr-2.5">
					<p>
						{{ stats?.title ?? 'No title' }}
					</p>
					<div class="block md:hidden">|</div>
					<p>
						{{ stats?.categoryName ?? 'No category' }}
					</p>
				</div>
				<IconEdit class="h-5 w-5 cursor-pointer" />
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</p>
				<p class="stats-display">
					{{ uptime }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</p>
				<p class="stats-display">
					{{ stats?.viewers ?? 0 }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</p>

				<p class="stats-display">
					{{ stats?.followers }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</p>
				<p class="stats-display">
					{{ stats?.chatMessages }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</p>
				<p class="stats-display">
					{{ stats?.subs }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</p>
				<p class="stats-display">
					{{ stats?.usedEmotes }}
				</p>
			</div>

			<div class="item stats-item md:block! hidden!">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</p>
				<p class="stats-display">
					{{ stats?.requestedSongs }}
				</p>
			</div>
		</div>

		<div class="ml-auto flex flex-wrap justify-end gap-2 flex-end items-center">
			<HeaderBotStatus />
			<HeaderProfile />
		</div>
	</div>

	<StreamInfoEditor
		v-model:open="streamInfoEditorOpen"
		:title="stats?.title"
		:category-id="stats?.categoryId"
		:category-name="stats?.categoryName"
	/>
</template>
