<script setup lang="ts">
import { NScrollbar } from 'naive-ui'
import { storeToRefs } from 'pinia'

import Pagination from '@/components/pagination.vue'
import Table from '@/components/table.vue'
import {
	useCommunityEmotesDetails,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details'
import {
	useCommunityEmotesDetailsUsersTop,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details-users-top'

const { isLoading, topPagination } = storeToRefs(useCommunityEmotesDetails())
const users = useCommunityEmotesDetailsUsersTop()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<Pagination
			:total="users.total"
			:table="users.table"
			:pagination="topPagination"
			@update:page="(page) => topPagination.pageIndex = page"
			@update:page-size="(pageSize) => topPagination.pageSize = pageSize"
		/>

		<NScrollbar style="max-height: 400px;" trigger="none">
			<Table :table="users.table" :is-loading="isLoading" hide-header>
				<template #empty-message>
					Empty
				</template>
			</Table>
		</NScrollbar>
	</div>
</template>
