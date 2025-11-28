<script setup lang="ts">
import { LoaderCircleIcon } from 'lucide-vue-next'
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

import { useNightbotIntegration } from './composables/use-nightbot-integration.js'

const nightbotIntegration = useNightbotIntegration()
const route = useRoute()

onMounted(async () => {
	const { code } = route.query
	if (typeof code !== 'string') {
		return window.close()
	}

	await nightbotIntegration.postCode.executeMutation({ input: { code } })
	nightbotIntegration.broadcastRefresh()
	window.close()
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<LoaderCircleIcon class="animate-spin size-12" />
	</div>
</template>
