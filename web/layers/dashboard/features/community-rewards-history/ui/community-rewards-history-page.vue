<script setup lang="ts">
import CommunityRewardsTable from './community-rewards-history-table.vue'

import CommunityRewardsHistorySearch
	from '~~/layers/dashboard/features/community-rewards-history/ui/community-rewards-history-search.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<CommunityRewardsHistorySearch />
		<slot name="pagination" />
		<CommunityRewardsTable />
		<slot name="pagination" />
	</div>
</template>
