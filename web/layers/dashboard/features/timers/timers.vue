<script lang="ts" setup>
import TimersCreateButton from '~~/layers/dashboard/features/timers/ui/timers-create-button.vue'
import TimersTable from '~~/layers/dashboard/features/timers/ui/timers-table.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('sidebar.timers') }}
		</template>

		<template #action>
			<TimersCreateButton />
		</template>

		<template #content>
			<TimersTable />
		</template>
	</PageLayout>
</template>
