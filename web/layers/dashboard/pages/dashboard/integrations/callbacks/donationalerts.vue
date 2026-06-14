<script setup lang="ts">
import { onMounted } from 'vue'
import { useIntegrations } from '~~/layers/dashboard/api/integrations/integrations.js'
import { useRoute } from 'vue-router'

definePageMeta({ layout: 'popup', middleware: 'auth' })

const integration = useIntegrations()
const executor = integration.donationAlertsPostCode()
const route = useRoute<'dashboard-integrations-callbacks-donationalerts'>()

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
		<Icon name="lucide:loader-circle" class="animate-spin size-12" />
	</div>
</template>
