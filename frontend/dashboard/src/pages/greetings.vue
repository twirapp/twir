<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import Greetings from '@/features/greetings/greetings.vue'
import GreetingsCreateButton from '@/features/greetings/ui/greetings-create-button.vue'
import PageLayout from '@/layout/page-layout.vue'

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
