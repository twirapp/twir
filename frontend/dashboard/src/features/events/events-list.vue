<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useEventsTable } from './composables/use-events-table'

import Table from '@/components/table.vue'
import { Input } from '@/components/ui/input'
import EventsCreateButton from '@/features/events/ui/events-create-button.vue'
import PageLayout from '@/layout/page-layout.vue'

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
				<Input v-model="eventsTable.search.value" placeholder="Search..." class="w-full lg:w-[40%]" />
				<Table :table="eventsTable.table" :is-loading="eventsTable.isLoading.value" />
			</div>
		</template>
	</PageLayout>
</template>
