<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { ChevronDownIcon } from 'lucide-vue-next'
import { NText } from 'naive-ui'
import { computed, h, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import StreamInfoEditor from '../stream-info-editor.vue'

import { useBotInfo, useBotJoinPart, useRealtimeDashboardStats } from '@/api'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { padTo2Digits } from '@/helpers/convertMillisToTime'

const { stats } = useRealtimeDashboardStats()
const stateMutation = useBotJoinPart()
const { data: botInfo, refetch: refetchBotStatus } = useBotInfo()

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

onMounted(async () => {
	refetchBotStatus()
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

async function changeChatState() {
	const action = botInfo.value?.enabled ? 'part' : 'join'

	await stateMutation.mutateAsync(action)
}
</script>

<template>
	<div class="flex flex-wrap justify-between bg-card w-full h-auto px-2 gap-2 border-b border-b-border min-h-12">
		<div class="flex flex-wrap gap-4 py-2">
			<div class="flex items-center cursor-pointer" @click="openInfoEditor">
				<div class="flex flex-col pr-2.5">
					<NText>
						{{ stats?.title ?? 'No title' }}
					</NText>
					<NText>
						{{ stats?.categoryName ?? 'No category' }}
					</NText>
				</div>
				<IconEdit class="h-5 w-5 cursor-pointer" />
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</NText>
				<NText class="stats-display">
					{{ uptime }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.viewers ?? 0 }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</NText>

				<NText class="stats-display">
					{{ stats?.followers }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.chatMessages }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.subs }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.usedEmotes }}
				</NText>
			</div>

			<div class="item stats-item">
				<NText :depth="3" class="stats-type">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</NText>
				<NText class="stats-display">
					{{ stats?.requestedSongs }}
				</NText>
			</div>
		</div>

		<div class="flex justify-end flex-end items-center">
			<Button v-if="!botInfo?.enabled" size="sm" @click="changeChatState">
				Join channel
			</Button>
			<DropdownMenu v-else>
				<DropdownMenuTrigger as-child>
					<Button variant="secondary" size="sm">
						<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="me-0.5 text-green-400 ping"><circle cx="12" cy="12" r="5" fill="currentColor"></circle><circle cx="12" cy="12" r="6" fill="currentColor" style="animation: 1s cubic-bezier(0, 0, 0.2, 1) 0s infinite normal none running ping; transform-origin: center center;"></circle></svg>
						{{ botInfo.botName }} online
						<ChevronDownIcon class="ml-2 size-4" />
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent>
					<DropdownMenuItem class="text-red-700" @click="changeChatState">
						Leave channel
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	</div>
</template>

<style scoped>
.stats-item {
	@apply flex flex-col justify-between min-w-8 rounded-md;
}

.stats-type {
	@apply text-xs;
}

.stats-display {
	@apply text-base tabular-nums;
}

@keyframes ping { 75%, 100% { transform: scale(2); opacity: 0; } }
</style>
