<script setup lang="ts">
import { watch } from 'vue';
import { RouterView, useRoute } from 'vue-router';
import { apiKeyRef } from '@/plugins/urql.ts';
import { loadEruda } from './helpers.js'

const route = useRoute()

// Set initial API key
if (route.params.apiKey) {
	apiKeyRef.value = route.params.apiKey as string
}

// Attach Eruda devtools
if (route.query.debug) {
	loadEruda()
}

// Watch for route changes to update API key
watch(
	() => route.params.apiKey,
	(newApiKey) => {
		if (newApiKey) {
			apiKeyRef.value = newApiKey as string
		}
	},
	{ immediate: true }
)

// if (import.meta.env.DEV) {
// 	document.body.style.backgroundColor = '#000';
// }
</script>

<template>
	<router-view />
</template>
