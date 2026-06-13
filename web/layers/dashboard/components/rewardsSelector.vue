<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useTwitchRewardsNew } from '@/api/twitch.ts'
import RewardFallbackImg from '@/assets/images/reward-fallback.png?url'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'

const props = defineProps<{
	multiple?: boolean
	clearable?: boolean
	onlyWithInput?: boolean
	placeholder?: string
}>()

const modelValue = defineModel<string | string[] | undefined | null>()

const { t } = useI18n()

const {
	data: rewardsData,
	fetching: isRewardsLoading,
	error: isRewardsError,
} = useTwitchRewardsNew()

interface RewardOption {
	id: string
	title: string
	image?: string
	color: string
	enabled: boolean
}

const rewardsOptions = computed<RewardOption[]>(() => {
	const rewards: RewardOption[] = []
	if (!rewardsData.value?.twitchRewards) return rewards

	for (const reward of rewardsData.value.twitchRewards) {
		if (props.onlyWithInput && !reward.userInputRequired) continue

		rewards.push({
			id: reward.id,
			title: reward.title,
			image: reward.imageUrls?.at(-1),
			color: reward.backgroundColor,
			enabled: reward.enabled,
		})
	}

	return rewards
})

const selectedReward = computed(() => {
	if (!modelValue.value || Array.isArray(modelValue.value)) return null
	return rewardsOptions.value.find(r => r.id === modelValue.value)
})
</script>

<template>
	<!-- Single select version -->
	<Select
		v-if="!multiple"
		v-model="modelValue"
		:disabled="isRewardsLoading || isRewardsError !== undefined"
	>
		<SelectTrigger class="w-full">
			<SelectValue :placeholder="placeholder ?? t('events.targetTwitchReward')">
				<div v-if="selectedReward" class="flex items-center gap-2">
					<Avatar class="h-5 w-5">
						<AvatarImage :src="selectedReward.image || RewardFallbackImg" />
						<AvatarFallback>
							<div
								class="w-full h-full"
								:style="{ backgroundColor: selectedReward.color }"
							/>
						</AvatarFallback>
					</Avatar>
					<span :class="cn(!selectedReward.enabled && 'text-muted-foreground')">
						{{ selectedReward.title }}
					</span>
				</div>
			</SelectValue>
		</SelectTrigger>
		<SelectContent>
			<SelectGroup>
				<SelectItem
					v-for="reward in rewardsOptions"
					:key="reward.id"
					:value="reward.id"
					:disabled="!reward.enabled"
				>
					<div class="flex items-center gap-2">
						<Avatar class="h-5 w-5">
							<AvatarImage :src="reward.image || RewardFallbackImg" />
							<AvatarFallback>
								<div
									class="w-full h-full"
									:style="{ backgroundColor: reward.color }"
								/>
							</AvatarFallback>
						</Avatar>
						<span :class="cn(!reward.enabled && 'text-muted-foreground')">
							{{ reward.title }}
						</span>
					</div>
				</SelectItem>
			</SelectGroup>
		</SelectContent>
	</Select>

	<!-- Multiple select version - fallback to native multi-select for now -->
	<div v-else class="space-y-2">
		<p class="text-sm text-muted-foreground">
			Multiple select for rewards is not yet implemented with shadcn components.
			Please use the single select version for now.
		</p>
	</div>
</template>
