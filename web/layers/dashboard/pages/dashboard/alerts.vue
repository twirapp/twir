<script setup lang="ts">
import Alerts from '~~/layers/dashboard/features/alerts/alerts.vue'
import AlertsCreateButton from '~~/layers/dashboard/features/alerts/ui/alerts-create-button.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('alerts.title') }}
		</template>

		<template #action>
			<AlertsCreateButton />
		</template>

		<template #content>
			<Alerts />
		</template>
	</PageLayout>
</template>
