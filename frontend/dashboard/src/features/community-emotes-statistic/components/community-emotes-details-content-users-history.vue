<script setup lang="ts">
import { NScrollbar } from 'naive-ui'

import Pagination from '@/components/pagination.vue'
import Table from '@/components/table.vue'
import {
	useCommunityEmotesDetails,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details'
import {
	useCommunityEmotesDetailsUsersHistory,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details-users-history'

const { isLoading, usagesPagination } = useCommunityEmotesDetails()
const { total, table } = useCommunityEmotesDetailsUsersHistory()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<Pagination
			:total="total"
			:table="table"
			:pagination="usagesPagination"
			@update:page="(page) => usagesPagination.pageIndex = page"
			@update:page-size="(pageSize) => usagesPagination.pageSize = pageSize"
		/>

		<NScrollbar style="max-height: 400px;" trigger="none">
			<Table :table="table" :is-loading="isLoading" hide-header>
				<template #empty-message>
					Empty
				</template>
			</Table>
		</NScrollbar>
	</div>
</template>
