<script setup lang="ts">
import { ref, watch } from 'vue'
import { ChevronsUpDown } from 'lucide-vue-next'

import { useBotJoinPart, useBotStatus } from '@/api'
import { BotJoinLeaveAction } from '@/gql/graphql.ts'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import CircleSvg from '@/assets/images/circle.svg?use'

const { botStatus, executeSubscription } = useBotStatus()
const stateMutation = useBotJoinPart()

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
	<Button
		v-if="!botStatus?.enabled"
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
				:disabled="waitingBotStatusData"
				class="flex items-center gap-0.5"
			>
				<CircleSvg class="circle text-green-400" />
				{{ botStatus?.botName ?? 'Bot' }} online
				<ChevronsUpDown class="ml-2 size-4" />
			</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent>
			<DropdownMenuItem class="text-red-700" @click="changeChatState">
				Leave channel
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>

<style>
@reference '@/assets/index.css';

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
