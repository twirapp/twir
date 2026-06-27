<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

import { computed } from 'vue'
import SecretsList from '~~/layers/dashboard/features/channels-secret/secrets.vue'
import VariablesList from '~~/layers/dashboard/features/variables/variables.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()

const tabs = computed<PageLayoutTab[]>(() => [
	{
		title: 'Variables',
		component: VariablesList,
		name: 'variables',
	},
	{
		title: 'Secrets',
		component: SecretsList,
		name: 'secrets',
	},
])
</script>

<template>
	<PageLayout
		:tabs="tabs"
		active-tab="variables"
	>
		<template #title>
			{{ t('sidebar.variables') }}
		</template>
	</PageLayout>
</template>
