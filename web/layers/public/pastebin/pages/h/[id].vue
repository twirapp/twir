<script setup lang="ts">
import HastebinLayout from '../components/HastebinLayout.vue'
import HastebinPage from '../components/HastebinPage.vue'

import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const route = useRoute()
const api = useOapi()

const { data, status, error } = await useAsyncData(
	'hastebin',
	async () => {
		const req = await api.v1.pastebinGetById(route.params.id as string)
		if (req.error) {
			throw req.error
		}

		return req.data
	},
)
const pasteStore = usePasteStore()

if (data.value) {
	pasteStore.setCurrentPaste(data.value)
}

const pageTitle = computed(() => {
	return data.value ? `Paste ${route.params.id}` : 'Loading paste...'
})
</script>

<template>
	<HastebinLayout :title="pageTitle">
		<div v-if="status === 'pending'" class="flex items-center justify-center w-full h-full text-white">
			<p>Loading paste...</p>
		</div>

		<div v-else-if="error" class="flex items-center justify-center w-full h-full text-white">
			<p>Error loading paste: {{ error.message }}</p>
		</div>

		<HastebinPage v-else :item="data" />
	</HastebinLayout>
</template>
