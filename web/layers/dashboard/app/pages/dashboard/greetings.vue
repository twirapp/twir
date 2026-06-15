<script setup lang="ts">
import Greetings from '~~/layers/dashboard/app/features/greetings/greetings.vue'
import GreetingsCreateButton from '~~/layers/dashboard/app/features/greetings/ui/greetings-create-button.vue'
import PageLayout from '~~/layers/dashboard/app/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('greetings.title') }}
		</template>

		<template #action>
			<GreetingsCreateButton />
		</template>

		<template #content>
			<Greetings />
		</template>
	</PageLayout>
</template>
