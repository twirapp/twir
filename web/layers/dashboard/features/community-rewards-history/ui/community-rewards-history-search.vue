<script setup lang="ts">
import { CheckIcon, ListFilterIcon } from 'lucide-vue-next'

import SearchBar from '#layers/dashboard/components/search-bar.vue'





import {
	useCommunityRewardsHistoryQuery,
} from '~/features/community-rewards-history/composables/community-rewards-history-query.ts'

const query = useCommunityRewardsHistoryQuery()
</script>

<template>
	<div class="flex gap-2">
		<SearchBar
			v-model="query.searchInput.value"
			placeholder="Search by username..."
		/>

		<UiPopover>
			<UiPopoverTrigger as-child>
				<UiButton variant="outline" size="sm" class="h-9">
					<ListFilterIcon class="mr-2 h-4 w-4" />
					Rewards filters

					<template v-if="query.query.value.rewardsIds?.length">
						<UiSeparator orientation="vertical" class="mx-2 h-4" />
						<UiBadge
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							Selected: {{ query.query.value.rewardsIds.length }}
						</UiBadge>
					</template>
				</UiButton>
			</UiPopoverTrigger>
			<UiPopoverContent class="w-fit-content p-0" align="end">
				<UiCommand>
					<UiCommandList>
						<UiCommandItem
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
								<CheckIcon class="h-4 w-4" />
							</div>
							<img v-if="reward.image" :src="reward.image" class="h-5 w-5 mr-2">
							<span>{{ reward.title }}</span>
						</UiCommandItem>
					</UiCommandList>
				</UiCommand>
			</UiPopoverContent>
		</UiPopover>
	</div>
</template>
