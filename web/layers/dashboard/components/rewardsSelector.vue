<script setup lang="ts">
import { computed, ref } from 'vue'


import { useTwitchRewardsNew } from '~~/layers/dashboard/api/twitch.js'
import RewardFallbackImg from '~~/layers/dashboard/assets/images/reward-fallback.png'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { cn } from '~~/layers/dashboard/lib/utils'

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

const open = ref(false)

const selectedIds = computed<string[]>(() => {
	if (!Array.isArray(modelValue.value)) return []
	return modelValue.value
})

const selectedIdsSet = computed(() => new Set(selectedIds.value))

function toggleReward(id: string) {
	const next = new Set(selectedIdsSet.value)
	if (next.has(id)) {
		next.delete(id)
	} else {
		next.add(id)
	}
	modelValue.value = Array.from(next)
}

function removeReward(id: string) {
	modelValue.value = selectedIds.value.filter(selectedId => selectedId !== id)
}

function clearRewards() {
	modelValue.value = []
}

function rewardById(id: string) {
	return rewardsOptions.value.find(r => r.id === id)
}
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

	<!-- Multiple select version -->
	<Popover v-else v-model:open="open">
		<PopoverTrigger as-child :disabled="isRewardsLoading || isRewardsError !== undefined">
			<Button
				variant="outline"
				role="combobox"
				:aria-expanded="open"
				class="w-full justify-between h-fit"
				:disabled="isRewardsLoading || isRewardsError !== undefined"
			>
				<div class="flex gap-1 flex-wrap items-center">
					<Badge
						v-for="id in selectedIds"
						:key="id"
						variant="default"
						class="mr-1 mb-1"
					>
						<template v-if="rewardById(id)">
							<Avatar class="h-4 w-4 mr-1">
								<AvatarImage :src="rewardById(id)!.image || RewardFallbackImg" />
								<AvatarFallback>
									<div
										class="w-full h-full"
										:style="{ backgroundColor: rewardById(id)!.color }"
									/>
								</AvatarFallback>
							</Avatar>
							{{ rewardById(id)!.title }}
						</template>
						<template v-else>
							{{ id }}
						</template>
						<button
							type="button"
							class="ml-1 ring-offset-background rounded-full outline-hidden focus:ring-2 focus:ring-ring focus:ring-offset-2"
							@click.stop="removeReward(id)"
						>
							<Icon name="lucide:x" class="h-3 w-3 text-background hover:text-background" />
						</button>
					</Badge>
					<span v-if="!selectedIds.length" class="text-muted-foreground font-normal">
						{{ placeholder ?? t('events.targetTwitchReward') }}
					</span>
				</div>
				<button
					v-if="clearable && selectedIds.length > 0"
					type="button"
					class="ml-2 shrink-0 text-muted-foreground hover:text-foreground"
					@click.stop.prevent="clearRewards"
				>
					<Icon name="lucide:x" class="h-4 w-4" />
				</button>
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-full p-0">
			<Command>
				<CommandInput placeholder="Search rewards..." />
				<CommandList>
					<CommandEmpty>No rewards found.</CommandEmpty>
					<CommandGroup>
						<CommandItem
							v-for="reward in rewardsOptions"
							:key="reward.id"
							:value="reward.id"
							:disabled="!reward.enabled"
							@select="toggleReward(reward.id)"
						>
							<span
								:class="cn(
									'mr-2 h-4 w-4 border border-primary rounded-sm flex items-center justify-center',
									selectedIdsSet.has(reward.id)
										? 'bg-primary text-primary-foreground'
										: 'opacity-50',
								)"
							>
								<Icon
									v-if="selectedIdsSet.has(reward.id)"
									name="lucide:check"
									class="size-4"
								/>
							</span>
							<Avatar class="h-4 w-4">
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
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
