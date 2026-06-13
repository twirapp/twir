<script setup lang="ts">
import Nightbot from '~/layers/dashboard/features/import/nightbot/nightbot-import.vue'
import StreamElements from '~/layers/dashboard/pages/import/streamelements.vue'
import Streamlabs from '~/layers/dashboard/pages/import/streamlabs.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<div class="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
		<nightbot />
		<stream-elements />
		<streamlabs />
	</div>
</template>
