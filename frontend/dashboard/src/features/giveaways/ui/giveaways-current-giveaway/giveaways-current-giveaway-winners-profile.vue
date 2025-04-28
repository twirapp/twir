<script setup lang="ts">
import { BubblesIcon, EyeIcon, GemIcon, HeartIcon, MessageSquareIcon, ShieldUserIcon, SmileIcon, SwordIcon } from 'lucide-vue-next'
import { computed, nextTick, toRef, watch } from 'vue'

import { useChannelUserInfo } from '@/api/users.ts'

const props = defineProps<{
	userId?: string
}>()

const userId = toRef(() => props.userId)

const {
	data: selectedWinnerChannelInformation,
	executeQuery: refetchWinnerChannelInfo,
} = useChannelUserInfo(userId, { manual: true })

watch(userId, async () => {
	await nextTick()
	await refetchWinnerChannelInfo({ requestPolicy: 'cache-and-network' })
}, {
	immediate: true,
})

const ONE_HOUR = 60 * 60 * 1000

const userHaveSomeRole = computed(() => {
	const info = selectedWinnerChannelInformation?.value?.channelUserInfo
	if (!info) return false

	return info.isMod || info.isVip || info.isSubscriber
})
</script>

<template>
	<div class="p-2 border-b border-border">
		<table class="table-auto w-full">
			<tbody>
				<tr>
					<td class="table-td">
						<MessageSquareIcon class="size-4" />
						Total messages
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.messages }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<EyeIcon class="size-4" />
						Watched time
					</td>
					<td>
						{{ `${(Number(selectedWinnerChannelInformation?.channelUserInfo.watchedMs) / ONE_HOUR).toFixed(1)}h` }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<SmileIcon class="size-4" />
						Used emotes
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.usedEmotes }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<BubblesIcon class="size-4" />
						Used channel points
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.usedChannelPoints }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<ShieldUserIcon class="size-4" />
						Roles
					</td>
					<td>
						<div class="flex gap-0.5">
							<span v-if="!userHaveSomeRole" class="font-light text-muted-foreground">
								No roles
							</span>

							<span v-if="selectedWinnerChannelInformation?.channelUserInfo.isMod" class="font-light">
								<SwordIcon class="size-4" />
								MOD
							</span>
							<span v-if="selectedWinnerChannelInformation?.channelUserInfo.isVip" class="font-light">
								<GemIcon class="size-4" />
								VIP
							</span>
							<span v-if="selectedWinnerChannelInformation?.channelUserInfo.isSubscriber" class="font-light">SUB</span>
						</div>
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<HeartIcon class="size-4" />
						Follower since
					</td>
					<td>
						<span v-if="selectedWinnerChannelInformation?.channelUserInfo.followerSince">
							{{ new Date(selectedWinnerChannelInformation.channelUserInfo.followerSince).toLocaleString() }}
						</span>
						<span v-else class="font-light text-muted-foreground">Not a follower</span>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped>
.table-td {
	@apply text-sm font-medium inline-flex items-center gap-2
}
</style>
