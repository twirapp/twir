<script setup lang="ts">
import { useChannelUserInfo } from '~~/layers/dashboard/api/users.js'

const props = defineProps<{
	userId?: string
}>()

const userId = toRef(() => props.userId)

const { data: selectedWinnerChannelInformation, executeQuery: refetchWinnerChannelInfo } =
	useChannelUserInfo(userId, { manual: true })

watch(
	userId,
	async () => {
		await nextTick()
		await refetchWinnerChannelInfo({ requestPolicy: 'cache-and-network' })
	},
	{
		immediate: true,
	}
)

const ONE_HOUR = 60 * 60 * 1000

const userHaveSomeRole = computed(() => {
	const info = selectedWinnerChannelInformation?.value?.channelUserInfo
	if (!info) return false

	return info.isMod || info.isVip || info.isSubscriber
})
</script>

<template>
	<div class="border-border border-b p-2">
		<table class="w-full table-auto">
			<tbody>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:message-square"
							class="size-4"
						/>
						Total messages
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.messages }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:eye"
							class="size-4"
						/>

						Watched time
					</td>
					<td>
						{{
							`${(Number(selectedWinnerChannelInformation?.channelUserInfo.watchedMs) / ONE_HOUR).toFixed(1)}h`
						}}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:smile"
							class="size-4"
						/>
						Used emotes
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.usedEmotes }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:bubbles"
							class="size-4"
						/>
						Used channel points
					</td>
					<td>
						{{ selectedWinnerChannelInformation?.channelUserInfo.usedChannelPoints }}
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:shield-user"
							class="size-4"
						/>
						Roles
					</td>
					<td>
						<div class="flex gap-0.5">
							<span
								v-if="!userHaveSomeRole"
								class="text-muted-foreground font-light"
							>
								No roles
							</span>

							<span
								v-if="selectedWinnerChannelInformation?.channelUserInfo.isMod"
								class="font-light"
							>
								MOD
							</span>
							<span
								v-if="selectedWinnerChannelInformation?.channelUserInfo.isVip"
								class="font-light"
							>
								VIP
							</span>
							<span
								v-if="selectedWinnerChannelInformation?.channelUserInfo.isSubscriber"
								class="font-light"
								>SUB</span
							>
						</div>
					</td>
				</tr>
				<tr>
					<td class="table-td">
						<Icon
							name="lucide:heart"
							class="size-4"
						/>
						Follower since
					</td>
					<td>
						<span v-if="selectedWinnerChannelInformation?.channelUserInfo.followerSince">
							{{
								new Date(
									selectedWinnerChannelInformation.channelUserInfo.followerSince
								).toLocaleString()
							}}
						</span>
						<span
							v-else
							class="text-muted-foreground font-light"
							>Not a follower</span
						>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped>
.table-td {
	display: inline-flex;
	align-items: center;
	gap: 0.5rem;
	font-size: 0.875rem;
	font-weight: 500;
	line-height: 1.25rem;
}
</style>
