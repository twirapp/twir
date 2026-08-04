<script setup lang="ts">
import Nightbot from '~~/layers/dashboard/features/import/nightbot/nightbot-import.vue'
import StreamElements from '~~/layers/dashboard/features/import/streamelements/streamelements-import.vue'
import Streamlabs from '~~/layers/dashboard/features/import/streamlabs/streamlabs-import.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<div class="flex flex-col gap-6 p-4 md:p-6">
		<header class="flex flex-col gap-2">
			<h1 class="text-2xl font-semibold tracking-tight">{{ $t('imports.title') }}</h1>
			<p class="max-w-3xl text-muted-foreground">{{ $t('imports.description') }}</p>
		</header>
		<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
			<Nightbot />
			<StreamElements />
			<Streamlabs />
		</div>
	</div>
</template>
