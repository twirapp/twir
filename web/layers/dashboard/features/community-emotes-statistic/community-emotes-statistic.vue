<script setup lang="ts">
import { useCommunityEmotesStatisticTable } from './composables/use-community-emotes-statistic-table'
import Table from './ui/community-emotes-table.vue'

import Pagination from '#layers/dashboard/components/pagination.vue'

const emotesTable = useCommunityEmotesStatisticTable()
</script>

<template>
	<Table>
		<template #pagination>
			<Pagination
				:total="emotesTable.totalEmotes.value"
				:table="emotesTable.table"
				:pagination="emotesTable.pagination.value"
				@update:page="(page) => emotesTable.pagination.value.pageIndex = page"
				@update:page-size="(pageSize) => emotesTable.pagination.value.pageSize = pageSize"
			/>
		</template>
	</Table>
</template>
