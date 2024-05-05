<script setup lang="ts">
import { NScrollbar } from 'naive-ui'
import { storeToRefs } from 'pinia'

import Pagination from '@/components/pagination.vue'
import Table from '@/components/table.vue'
import {
	useCommunityEmotesDetails,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details'
import {
	useCommunityEmotesDetailsUsers,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details-users'

const { isLoading } = storeToRefs(useCommunityEmotesDetails())
const users = useCommunityEmotesDetailsUsers()
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<Pagination
			:total="users.total"
			:table="users.table"
			:pagination="users.pagination"
			@update:page="(page) => users.pagination.pageIndex = page"
			@update:page-size="(pageSize) => users.pagination.pageSize = pageSize"
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
