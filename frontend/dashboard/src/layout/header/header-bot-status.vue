<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ChevronsUpDown, Loader2, LogIn, LogOut } from 'lucide-vue-next'

import { useBotJoinPart, useBotStatuses } from '@/api/dashboard'
import { BotJoinLeaveAction } from '@/gql/graphql.ts'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import KickIcon from '@/components/kick-icon.vue'
import CircleSvg from '@/assets/images/circle.svg?use'

const { botStatuses, executeSubscription } = useBotStatuses()
const stateMutation = useBotJoinPart()

const pendingStatusKeys = ref<Set<string>>(new Set())
const hasReceivedStatus = ref(false)

const sortedBotStatuses = computed(() => {
	return [...botStatuses.value].sort((a, b) => {
		if (a.platform === b.platform) return a.botName.localeCompare(b.botName)
		if (a.platform === 'twitch') return -1
		if (b.platform === 'twitch') return 1
		return a.platform.localeCompare(b.platform)
	})
})

const enabledStatusesCount = computed(() => {
	return sortedBotStatuses.value.filter((status) => status.enabled).length
})

const allStatusesEnabled = computed(() => {
	return sortedBotStatuses.value.length > 0 && enabledStatusesCount.value === sortedBotStatuses.value.length
})

const statusSummary = computed(() => {
	if (!hasReceivedStatus.value) return 'Bot status'
	if (sortedBotStatuses.value.length === 0) return 'Bot offline'

	if (sortedBotStatuses.value.length === 1) {
		const status = sortedBotStatuses.value[0]
		return `${formatPlatformName(status.platform)} ${status.enabled ? 'online' : 'disabled'}`
	}

	if (enabledStatusesCount.value === 0) return 'All bots disabled'
	if (allStatusesEnabled.value) return `${sortedBotStatuses.value.length} platforms online`
	return `${enabledStatusesCount.value}/${sortedBotStatuses.value.length} platforms online`
})

watch(botStatuses, () => {
	hasReceivedStatus.value = true
	pendingStatusKeys.value = new Set()
})

function statusKey(status: { dashboardId: string; platform: string }) {
	return `${status.dashboardId}:${status.platform}`
}

function formatPlatformName(platform: string) {
	if (platform === 'kick') return 'Kick'
	if (platform === 'twitch') return 'Twitch'
	return platform || 'Bot'
}

function isStatusPending(status: { dashboardId: string; platform: string }) {
	return pendingStatusKeys.value.has(statusKey(status))
}

async function changeChatState(status: { dashboardId: string; platform: string; channelName: string; enabled: boolean }) {
	const key = statusKey(status)
	if (pendingStatusKeys.value.has(key)) {
		return
	}

	const nextPending = new Set(pendingStatusKeys.value)
	nextPending.add(key)
	pendingStatusKeys.value = nextPending

	const action = status.enabled ? BotJoinLeaveAction.Leave : BotJoinLeaveAction.Join
	const result = await stateMutation.executeMutation({
		action,
		dashboardId: status.dashboardId,
		platform: status.platform,
	})

	if (result.error) {
		const pendingAfterError = new Set(pendingStatusKeys.value)
		pendingAfterError.delete(key)
		pendingStatusKeys.value = pendingAfterError
		return
	}

	executeSubscription()
}
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger as-child>
			<Button
				variant="secondary"
				:disabled="!hasReceivedStatus"
				class="flex items-center gap-2"
			>
				<CircleSvg
					class="circle"
					:class="allStatusesEnabled ? 'text-green-400' : 'text-red-400'"
				/>
				<div class="flex items-center gap-1">
					<template v-for="status in sortedBotStatuses" :key="statusKey(status)">
						<KickIcon v-if="status.platform === 'kick'" class="size-4 text-[#53FC18]" />
						<Badge
							v-else-if="status.platform === 'twitch'"
							variant="outline"
							class="h-4 px-1 text-[10px]"
						>
							T
						</Badge>
					</template>
				</div>
				<span class="max-w-44 truncate">{{ statusSummary }}</span>
				<ChevronsUpDown class="size-4" />
			</Button>
		</DropdownMenuTrigger>
			<DropdownMenuContent align="end" class="w-72">
				<DropdownMenuLabel>Bot platforms</DropdownMenuLabel>
				<DropdownMenuSeparator />
				<DropdownMenuItem
					v-for="status in sortedBotStatuses"
					:key="statusKey(status)"
					class="flex items-center gap-3"
					:disabled="isStatusPending(status)"
					@select.prevent
					@click="changeChatState(status)"
				>
				<div class="flex size-7 items-center justify-center rounded-md border border-border bg-background">
					<KickIcon v-if="status.platform === 'kick'" class="size-4 text-[#53FC18]" />
					<Badge
						v-else-if="status.platform === 'twitch'"
						variant="outline"
						class="h-4 px-1 text-[10px]"
					>
						T
					</Badge>
				</div>
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm font-medium">
						{{ formatPlatformName(status.platform) }}
						<span class="text-muted-foreground">{{ status.channelName || status.botName || 'Bot' }}</span>
					</p>
					<p class="flex items-center gap-1 text-xs text-muted-foreground">
						<CircleSvg
							class="size-3"
							:class="status.enabled ? 'text-green-400' : 'text-red-400'"
						/>
						{{ status.enabled ? 'Online' : 'Disabled' }}
						<span v-if="status.botName" class="truncate">via {{ status.botName }}</span>
					</p>
				</div>
					<Loader2 v-if="isStatusPending(status)" class="size-4 animate-spin text-muted-foreground" />
					<LogOut v-else-if="status.enabled" class="size-4 text-red-500" />
					<LogIn v-else class="size-4 text-green-500" />
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
</template>

<style>
@reference '@/assets/index.css';

.circle {
	@apply size-4;
}
</style>
