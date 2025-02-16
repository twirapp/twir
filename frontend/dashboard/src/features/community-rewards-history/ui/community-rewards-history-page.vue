<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'

import CommunityRewardsTable from './community-rewards-history-table.vue'

import {
	useCommunityRewardsTable,
} from '@/features/community-rewards-history/composables/community-rewards-history-table.ts'
import CommunityRewardsHistorySearch
	from '@/features/community-rewards-history/ui/community-rewards-history-search.vue'

const table = useCommunityRewardsTable()

useIntervalFn(() => {
	table.refresh()
}, 1000)
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<CommunityRewardsHistorySearch />
		<slot name="pagination" />
		<CommunityRewardsTable />
		<slot name="pagination" />
	</div>
</template>
