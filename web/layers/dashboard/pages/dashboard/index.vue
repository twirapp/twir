<script setup lang="ts">
definePageMeta({
	layout: 'default',
	noPadding: true,
})

const authStore = useDashboardAuth()
const { t } = useI18n()
</script>

<template>
	<div class="p-6">
		<h1 class="text-3xl font-bold mb-4">
			{{ t('dashboard.title') }}
		</h1>

		<div v-if="authStore.user" class="space-y-4">
			<p class="text-lg">
				{{ t('dashboard.welcome') }}, {{ authStore.user.displayName }}!
			</p>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				<div class="p-4 border rounded-lg">
					<h3 class="font-semibold">{{ t('dashboard.stats.title') }}</h3>
					<p class="text-sm text-muted-foreground">{{ t('dashboard.stats.description') }}</p>
				</div>
			</div>
		</div>
	</div>
</template>
