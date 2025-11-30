<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { LoaderCircleIcon } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import { useValorantIntegration } from '@/features/integrations/composables/valorant/use-valorant-integration.ts'

const router = useRouter()
const { broadcastRefresh, postCode } = useValorantIntegration()

const errorResult = ref<string>()

onMounted(async () => {
	const { code, error } = router.currentRoute.value.query
	if (!code || typeof code != 'string') {
		errorResult.value = 'Something unexpected happend, contact developers in disocrd to get help'
		return
	}

	if (error) {
		errorResult.value = `Cannot connect valorant due error from riot: ${error}`
		return
	}

	const resultError = await postCode(code)
	if (resultError) {
		errorResult.value = `Cannot connect valorant: ${resultError}`
		return
	}

	broadcastRefresh()
	window.close()
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<LoaderCircleIcon v-if="!errorResult" class="animate-spin size-12" />
		<div v-else class="p-4 bg-red-950/50 text-2xl">
			{{ errorResult }}
		</div>
	</div>
</template>
