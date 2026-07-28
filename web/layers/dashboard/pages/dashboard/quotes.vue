<script setup lang="ts">
import Quotes from '~~/layers/dashboard/features/quotes/quotes.vue'
import QuotesCreateButton from '~~/layers/dashboard/features/quotes/ui/quotes-create-button.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('quotes.title') }}
		</template>
		<template #action>
			<QuotesCreateButton />
		</template>
		<template #content>
			<Quotes />
		</template>
	</PageLayout>
</template>
