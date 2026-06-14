<script setup lang="ts">

import { useEventsTable } from './composables/use-events-table'

import Table from '~~/layers/dashboard/components/table.vue'
import { Input } from '@/components/ui/input'
import EventsCreateButton from '~~/layers/dashboard/features/events/ui/events-create-button.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'
import EventsListSelectType from '~~/layers/dashboard/features/events/ui/events-list-select-type.vue'

const eventsTable = useEventsTable()
const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('sidebar.events') }}
		</template>

		<template #action>
			<EventsCreateButton />
		</template>

		<template #content>
			<div class="flex flex-col w-full gap-4">
				<div class="flex items-center flex-wrap justify-between gap-4">
					<Input
						v-model="eventsTable.search.value"
						placeholder="Search..."
						class="w-full lg:w-[40%]"
					/>
					<EventsListSelectType />
				</div>
				<Table :table="eventsTable.table" :is-loading="eventsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>
