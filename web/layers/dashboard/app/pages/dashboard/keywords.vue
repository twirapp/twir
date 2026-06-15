<script setup lang="ts">
import Keywords from '~~/layers/dashboard/app/features/keywords/keywords.vue'
import KeywordsCreateButton from '~~/layers/dashboard/app/features/keywords/ui/keywords-create-button.vue'
import PageLayout from '~~/layers/dashboard/app/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('keywords.title') }}
		</template>

		<template #action>
			<KeywordsCreateButton />
		</template>

		<template #content>
			<Keywords />
		</template>
	</PageLayout>
</template>
