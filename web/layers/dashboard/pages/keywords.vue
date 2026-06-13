<script setup lang="ts">
import Keywords from '~~/layers/dashboard/features/keywords/keywords.vue'
import KeywordsCreateButton from '~~/layers/dashboard/features/keywords/ui/keywords-create-button.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

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
