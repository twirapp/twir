<script lang="ts" setup>
import { useI18n } from 'vue-i18n'

import TimersCreateButton from '@/features/timers/ui/timers-create-button.vue'
import TimersTable from '@/features/timers/ui/timers-table.vue'
import PageLayout from '@/layout/page-layout.vue'

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
