<script setup lang="ts">
import { IconEdit } from '@tabler/icons-vue'
import { useIntervalFn } from '@vueuse/core'
import { intervalToDuration } from 'date-fns'
import { ChevronDownIcon } from 'lucide-vue-next'
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import StreamInfoEditor from '../stream-info-editor.vue'

import CircleSvg from '@/assets/images/circle.svg?use'
import { useRealtimeDashboardStats } from '@/api'
import { useBotJoinPart, useBotStatus } from '@/api/dashboard'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { BotJoinLeaveAction } from '@/gql/graphql.ts'
import { padTo2Digits } from '@/helpers/convertMillisToTime'

const { stats } = useRealtimeDashboardStats()
const stateMutation = useBotJoinPart()
const { botStatus, executeSubscription } = useBotStatus()

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

const waitingBotStatusData = ref(true)

watch(botStatus, () => {
	waitingBotStatusData.value = false
})

async function changeChatState() {
	const action = botStatus.value?.enabled ? BotJoinLeaveAction.Leave : BotJoinLeaveAction.Join

	waitingBotStatusData.value = true
	await stateMutation.executeMutation({ action })
	executeSubscription()
}
</script>

<template>
	<div
		class="flex flex-wrap justify-between bg-card w-full h-auto px-2 gap-2 border-b border-b-border min-h-12"
	>
		<div class="flex flex-wrap gap-4 py-2">
			<div class="flex items-center cursor-pointer" @click="openInfoEditor">
				<div class="flex flex-col pr-2.5">
					<p>
						{{ stats?.title ?? 'No title' }}
					</p>
					<p>
						{{ stats?.categoryName ?? 'No category' }}
					</p>
				</div>
				<IconEdit class="h-5 w-5 cursor-pointer" />
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.uptime`) }}
				</p>
				<p class="stats-display">
					{{ uptime }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.viewers`) }}
				</p>
				<p class="stats-display">
					{{ stats?.viewers ?? 0 }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.followers`) }}
				</p>

				<p class="stats-display">
					{{ stats?.followers }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.messages`) }}
				</p>
				<p class="stats-display">
					{{ stats?.chatMessages }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.subs`) }}
				</p>
				<p class="stats-display">
					{{ stats?.subs }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.usedEmotes`) }}
				</p>
				<p class="stats-display">
					{{ stats?.usedEmotes }}
				</p>
			</div>

			<div class="item stats-item">
				<p class="stats-type">
					{{ t(`dashboard.statsWidgets.requestedSongs`) }}
				</p>
				<p class="stats-display">
					{{ stats?.requestedSongs }}
				</p>
			</div>
		</div>

		<div class="flex justify-end flex-end items-center">
			<Button
				v-if="!botStatus?.enabled"
				size="sm"
				:disabled="waitingBotStatusData"
				@click="changeChatState"
				class="flex items-center gap-0.5"
				variant="secondary"
			>
				<CircleSvg class="circle text-red-400" />
				{{ botStatus?.botName ?? 'Bot' }} disabled, click to join channel
			</Button>
			<DropdownMenu v-else>
				<DropdownMenuTrigger as-child>
					<Button
						variant="secondary"
						size="sm"
						:disabled="waitingBotStatusData"
						class="flex items-center gap-0.5"
					>
						<CircleSvg class="circle text-green-400" />
						{{ botStatus?.botName ?? 'Bot' }} online
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

	<StreamInfoEditor
		v-model:open="streamInfoEditorOpen"
		:title="stats?.title"
		:category-id="stats?.categoryId"
		:category-name="stats?.categoryName"
	/>
</template>

<style>
.stats-item {
	@apply flex flex-col justify-between min-w-8 rounded-md;
}

.stats-type {
	@apply text-xs;
}

.stats-display {
	@apply text-base tabular-nums;
}

.circle {
	@apply size-6;
}

@keyframes ping {
	75%,
	100% {
		transform: scale(2);
		opacity: 0;
	}
}
</style>
