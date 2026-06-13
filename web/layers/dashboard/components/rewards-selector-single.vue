<script setup lang="ts">
import { CheckIcon, XIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import type { GetChannelRewardsQuery } from '@/gql/graphql.ts'

import { useTwitchRewardsNew } from '@/api/twitch.ts'
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
		<Popover v-model:open="open">
			<PopoverTrigger as-child>
				<Button
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
				</Button>
			</PopoverTrigger>
			<PopoverContent class="p-0" align="center">
				<Command :filter-function="filterFunction">
					<CommandInput placeholder="Change status..." />
					<CommandList>
						<CommandEmpty>No results found.</CommandEmpty>
						<CommandGroup>
							<CommandItem
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
							</CommandItem>
						</CommandGroup>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
		<Button v-if="deselect" size="icon" variant="secondary" @click="handleDeselect">
			<XIcon />
		</Button>
	</div>
</template>
