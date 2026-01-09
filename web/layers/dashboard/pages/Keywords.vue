<script setup lang="ts">


import Keywords from '~/features/keywords/keywords.vue'
import KeywordsCreateButton from '~/features/keywords/ui/keywords-create-button.vue'
import PageLayout from '~/layout/page-layout.vue'

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
