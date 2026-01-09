<script setup lang="ts">
import { useTwitchRewardsNew } from '#layers/dashboard/api/twitch.ts'
import RewardFallbackImg from '#layers/dashboard/assets/images/reward-fallback.png?url'
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
	<UiSelect
		v-if="!multiple"
		v-model="modelValue"
		:disabled="isRewardsLoading || isRewardsError !== undefined"
	>
		<UiSelectTrigger class="w-full">
			<UiSelectValue :placeholder="placeholder ?? t('events.targetTwitchReward')">
				<div v-if="selectedReward" class="flex items-center gap-2">
					<UiAvatar class="h-5 w-5">
						<UiAvatarImage :src="selectedReward.image || RewardFallbackImg" />
						<UiAvatarFallback>
							<div
								class="w-full h-full"
								:style="{ backgroundColor: selectedReward.color }"
							/>
						</UiAvatarFallback>
					</UiAvatar>
					<span :class="cn(!selectedReward.enabled && 'text-muted-foreground')">
						{{ selectedReward.title }}
					</span>
				</div>
			</UiSelectValue>
		</UiSelectTrigger>
		<UiSelectContent>
			<UiSelectGroup>
				<UiSelectItem
					v-for="reward in rewardsOptions"
					:key="reward.id"
					:value="reward.id"
					:disabled="!reward.enabled"
				>
					<div class="flex items-center gap-2">
						<UiAvatar class="h-5 w-5">
							<UiAvatarImage :src="reward.image || RewardFallbackImg" />
							<UiAvatarFallback>
								<div
									class="w-full h-full"
									:style="{ backgroundColor: reward.color }"
								/>
							</UiAvatarFallback>
						</UiAvatar>
						<span :class="cn(!reward.enabled && 'text-muted-foreground')">
							{{ reward.title }}
						</span>
					</div>
				</UiSelectItem>
			</UiSelectGroup>
		</UiSelectContent>
	</UiSelect>

	<!-- Multiple select version - fallback to native multi-select for now -->
	<div v-else class="space-y-2">
		<p class="text-sm text-muted-foreground">
			Multiple select for rewards is not yet implemented with shadcn components.
			Please use the single select version for now.
		</p>
	</div>
</template>
