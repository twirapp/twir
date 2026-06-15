<script setup lang="ts">
import SearchBar from '~~/layers/dashboard/components/search-bar.vue'
import Table from '~~/layers/dashboard/components/table.vue'

import { useCommunityEmotesStatisticFilters } from '../composables/use-community-emotes-statistic-filters'
import { useCommunityEmotesStatisticTable } from '../composables/use-community-emotes-statistic-table'

const { t } = useI18n()
const emotesStatisticTable = useCommunityEmotesStatisticTable()
const emotesStatisticFilter = useCommunityEmotesStatisticFilters()
</script>

<template>
	<div class="flex w-full flex-col gap-4">
		<SearchBar v-model="emotesStatisticFilter.searchInput.value" />
		<slot name="pagination" />
		<Table
			:table="emotesStatisticTable.table"
			:is-loading="emotesStatisticTable.isLoading.value"
		>
			<template #empty-message>
				{{ t('community.users.table.empty') }}
			</template>
		</Table>
		<slot name="pagination" />
	</div>
</template>
