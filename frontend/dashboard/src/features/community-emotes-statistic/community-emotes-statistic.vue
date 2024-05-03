<script setup lang="ts">
import Table from './components/community-emotes-table.vue'
import { useCommunityEmotesStatisticTable } from './composables/use-community-emotes-statistic-table'

import Pagination from '@/components/pagination.vue'

const emotesTable = useCommunityEmotesStatisticTable()
</script>

<template>
	<Table>
		<template #pagination>
			<Pagination
				:total="emotesTable.totalEmotes"
				:table="emotesTable.table"
				:pagination="emotesTable.pagination"
				@update:page="(page) => emotesTable.pagination.pageIndex = page"
				@update:page-size="(pageSize) => emotesTable.pagination.pageSize = pageSize"
			/>
		</template>
	</Table>
</template>
