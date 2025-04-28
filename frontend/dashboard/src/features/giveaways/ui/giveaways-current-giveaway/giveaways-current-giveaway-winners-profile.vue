<script setup lang="ts">
import { BubblesIcon, EyeIcon, HeartIcon, MessageSquareIcon, ShieldUserIcon, SmileIcon } from 'lucide-vue-next'
import { watch } from 'vue'

import { useChannelUserInfo } from '@/api/users.ts'

const props = defineProps<{
	userId?: string
}>()

const {
	data: selectedWinnerChannelInformation,
	executeQuery: refetchWinnerChannelInfo,
} = useChannelUserInfo(props.userId)

watch(() => props.userId, () => {
	refetchWinnerChannelInfo()
})
</script>

<template>
	<div v-if="selectedWinnerChannelInformation" class="shrink-0 p-2 border-b border-border flex flex-col space-y-1">
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<MessageSquareIcon class="size-4" />
				Total messages
			</h3>
			<div class="ml-auto">
				{{ selectedWinnerChannelInformation.channelUserInfo.messages }}
			</div>
		</div>
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<EyeIcon class="size-4" />
				Watched time
			</h3>
			<div class="ml-auto">
				{{ selectedWinnerChannelInformation.channelUserInfo.watchedMs }}
			</div>
		</div>
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<SmileIcon class="size-4" />
				Used emotes
			</h3>
			<div class="ml-auto">
				{{ selectedWinnerChannelInformation.channelUserInfo.usedEmotes }}
			</div>
		</div>
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<BubblesIcon class="size-4" />
				Used channel points
			</h3>
			<div class="ml-auto">
				{{ selectedWinnerChannelInformation.channelUserInfo.usedChannelPoints }}
			</div>
		</div>
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<ShieldUserIcon class="size-4" />
				Roles
			</h3>
			<div class="ml-auto flex gap-0.5">
				<span v-if="selectedWinnerChannelInformation.channelUserInfo.isMod" class="font-light">MOD</span>
				<span v-if="selectedWinnerChannelInformation.channelUserInfo.isVip" class="font-light">VIP</span>
				<span v-if="selectedWinnerChannelInformation.channelUserInfo.isSubscriber" class="font-light">SUB</span>
			</div>
		</div>
		<div class="flex">
			<h3 class="text-sm font-medium flex items-center gap-2">
				<HeartIcon class="size-4" />
				Follower since
			</h3>
			<div class="ml-auto">
				<span v-if="selectedWinnerChannelInformation.channelUserInfo.followerSince">
					{{ new Date(selectedWinnerChannelInformation.channelUserInfo.followerSince).toLocaleString() }}
				</span>
				<span v-else class="font-light text-muted-foreground">Not a follower</span>
			</div>
		</div>
	</div>
</template>
