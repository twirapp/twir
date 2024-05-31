<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import Alerts from '@/features/alerts/alerts.vue'
import AlertsCreateButton from '@/features/alerts/ui/alerts-create-button.vue'
import PageLayout from '@/layout/page-layout.vue'

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
