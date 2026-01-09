<script setup lang="ts">
import { LoaderCircleIcon } from 'lucide-vue-next'
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useIntegrations } from '#layers/dashboard/api/integrations/integrations.ts'

const integration = useIntegrations()
const executor = integration.donationAlertsPostCode()
const route = useRoute()

onMounted(async () => {
	const { code } = route.query
	if (typeof code !== 'string') {
		return window.close()
	}

	await executor.executeMutation({ code })
	integration.broadcastRefresh()
	window.close()
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<LoaderCircleIcon class="animate-spin size-12" />
	</div>
</template>
