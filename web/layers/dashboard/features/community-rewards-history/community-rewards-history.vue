<script setup lang="ts">
import { useCommunityRewardsTable } from './composables/community-rewards-history-table'
import CommunityRewardsPage from './ui/community-rewards-history-page.vue'

import Pagination from '@/components/pagination.vue'
import {
	useCommunityRewardsHistoryQuery,
} from '@/features/community-rewards-history/composables/community-rewards-history-query.ts'

const rewardsTable = useCommunityRewardsTable()
const query = useCommunityRewardsHistoryQuery()
</script>

<template>
	<CommunityRewardsPage>
		<template #pagination>
			<Pagination
				:total="rewardsTable.total.value"
				:table="rewardsTable.table"
				:pagination="query.pagination.value"
				@update:page="(page) => query.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => query.pagination.value.pageSize = pageSize"
			/>
		</template>
	</CommunityRewardsPage>
</template>
