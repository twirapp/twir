<script setup lang="ts">
import { CheckIcon, XIcon } from 'lucide-vue-next'

import type { GetChannelRewardsQuery } from '@/gql/graphql.ts'

import { useTwitchRewardsNew } from '#layers/dashboard/api/twitch.ts'

defineProps<{
	deselect?: boolean
	requireInput?: boolean
}>()

const open = ref(false)
const selectedReward = ref<GetChannelRewardsQuery['twitchRewards'][number] | undefined>()
const modelValue = defineModel<string | undefined | null>()

const {
	data: rewardsData,
} = useTwitchRewardsNew()

watch(modelValue, (value) => {
	if (value === undefined) {
		selectedReward.value = undefined
		return
	}

	selectedReward.value = rewardsData.value?.twitchRewards.find((reward) => reward.id === value)
})

function handleDeselect() {
	modelValue.value = null
}

const searchTerm = ref('')

function filterFunction(_items: any, searchInput: string) {
	searchTerm.value = searchInput
	return _items
}

const filteredRewards = computed(() => {
	if (!rewardsData.value?.twitchRewards) return []

	return rewardsData.value.twitchRewards.filter((reward) => reward.title.toLowerCase().includes(searchTerm.value.toLowerCase()))
})
</script>

<template>
	<div class="flex flex-row gap-2 flex-wrap">
		<UiPopover v-model:open="open">
			<UiPopoverTrigger as-child>
				<UiButton
					variant="outline"
					size="sm"
					class="flex-1 h-auto"
				>
					<div v-if="selectedReward" class="flex flex-row gap-2 items-center">
						<img :src="selectedReward?.imageUrls?.at(-1)" class="size-6" />
						<span>
							{{ selectedReward.title }}
						</span>
					</div>
					<template v-else>
						Select reward
					</template>
				</UiButton>
			</UiPopoverTrigger>
			<UiPopoverContent class="p-0" align="center">
				<UiCommand :filter-function="filterFunction">
					<UiCommandInput placeholder="Change status..." />
					<UiCommandList>
						<UiCommandEmpty>No results found.</UiCommandEmpty>
						<UiCommandGroup>
							<UiCommandItem
								v-for="reward in filteredRewards"
								:key="reward.id"
								:value="reward.id"
								class="mt-0.5"
								:class="modelValue === reward.id && 'bg-primary/10'"
								:disabled="!reward.userInputRequired && requireInput"
								@select="() => {
									modelValue = reward.id
									open = false
								}"
							>
								<div class="flex items-center gap-2 flex-row w-full">
									<img :src="reward.imageUrls?.at(-1)" class="size-6" />
									<span>
										{{ reward.title }}
									</span>
									<CheckIcon v-if="modelValue === reward.id" class="ml-auto" />
								</div>
							</UiCommandItem>
						</UiCommandGroup>
					</UiCommandList>
				</UiCommand>
			</UiPopoverContent>
		</UiPopover>
		<UiButton v-if="deselect" size="icon" variant="secondary" @click="handleDeselect">
			<XIcon />
		</UiButton>
	</div>
</template>
