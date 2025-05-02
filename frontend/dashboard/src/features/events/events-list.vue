<script setup lang="ts">
import { useEventsTable } from './composables/use-events-table'
import EventsInformation from './ui/events-information.vue'

import Table from '@/components/table.vue'
import PageLayout from '@/layout/page-layout.vue'

const eventsTable = useEventsTable()
</script>

<template>
	<PageLayout>
		<template #title>
			Events
		</template>

		<template #content>
			<div class="flex flex-col w-full gap-4">
				<EventsInformation />
				<Table :table="eventsTable.table" :is-loading="eventsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>
