<script setup lang="ts">
definePageMeta({
	layout: 'default',
})

const { t } = useI18n()
const router = useRouter()
</script>

<template>
	<div class="flex flex-col items-center justify-center h-full">
		<h1 class="text-4xl font-bold mb-4">
			{{ t('errors.forbidden.title') }}
		</h1>
		<p class="text-lg text-muted-foreground mb-8">
			{{ t('errors.forbidden.description') }}
		</p>
		<UiButton @click="router.push('/dashboard')">
			{{ t('errors.forbidden.back') }}
		</UiButton>
	</div>
</template>
