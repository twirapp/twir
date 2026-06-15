<script setup lang="ts">
import SearchBar from '~~/layers/dashboard/app/components/search-bar.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'
import {
	useCommunityRewardsHistoryQuery,
} from '~~/layers/dashboard/app/features/community-rewards-history/composables/community-rewards-history-query.js'

const query = useCommunityRewardsHistoryQuery()
</script>

<template>
	<div class="flex gap-2">
		<SearchBar
			v-model="query.searchInput.value"
			placeholder="Search by username..."
		/>

		<Popover>
			<PopoverTrigger as-child>
				<Button variant="outline" size="sm" class="h-9">
					<Icon name="lucide:list-filter" class="mr-2 h-4 w-4" />
					Rewards filters

					<template v-if="query.query.value.rewardsIds?.length">
						<Separator orientation="vertical" class="mx-2 h-4" />
						<Badge
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							Selected: {{ query.query.value.rewardsIds.length }}
						</Badge>
					</template>
				</Button>
			</PopoverTrigger>
			<PopoverContent class="w-fit-content p-0" align="end">
				<Command>
					<CommandList>
						<CommandItem
							v-for="reward of query.rewardsOptions.value"
							:key="reward.id"
							:value="reward.id"
							@select="query.handleRewardFilter(reward.id)"
						>
							<div
								class="mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary"
								:class="[query.query.value.rewardsIds?.includes(reward.id)
									? 'bg-primary text-primary-foreground'
									: 'opacity-50 [&_svg]:invisible',
								]"
							>
								<Icon name="lucide:check" class="h-4 w-4" />
							</div>
							<img v-if="reward.image" :src="reward.image" class="h-5 w-5 mr-2">
							<span>{{ reward.title }}</span>
						</CommandItem>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	</div>
</template>
