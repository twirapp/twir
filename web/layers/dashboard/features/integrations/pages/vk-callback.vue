<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useVKIntegration } from '~~/layers/dashboard/features/integrations/composables/vk/use-vk-integration.js'

const router = useRouter()
const { broadcastRefresh, postCode } = useVKIntegration()

const errorResult = ref<string>()

onMounted(async () => {
	const { code, error } = router.currentRoute.value.query
	if (!code || typeof code != 'string') {
		errorResult.value = 'Something unexpected happend, contact developers in discord to get help'
		return
	}

	if (error) {
		errorResult.value = `Cannot connect VK due error: ${error}`
		return
	}

	const resultError = await postCode(code)
	if (resultError) {
		errorResult.value = `Cannot connect VK: ${resultError}`
		return
	}

	broadcastRefresh()
	window.close()
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<Icon name="lucide:loader-circle" v-if="!errorResult" class="animate-spin size-12" />
		<div v-else class="p-4 bg-red-950/50 text-2xl">
			{{ errorResult }}
		</div>
	</div>
</template>
