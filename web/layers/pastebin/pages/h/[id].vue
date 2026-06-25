<script setup lang="ts">
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'
import { useRoute } from 'vue-router'

import HastebinLayout from '../../components/HastebinLayout.vue'
import HastebinPage from '../../components/HastebinPage.vue'

const route = useRoute<'h-id'>()
const api = useOapi()

const { data, status, error } = await useAsyncData('hastebin', async () => {
	const req = await api.v1.pastebinGetById(route.params.id)
	if (req.error) {
		throw req.error
	}

	return req.data
})

const pasteStore = usePasteStore()

if (data.value?.data) {
	pasteStore.setCurrentPaste(data.value?.data)
}

const pageTitle = computed(() => {
	return data.value ? `Paste ${route.params.id}` : 'Loading paste...'
})
</script>

<template>
	<HastebinLayout :title="pageTitle">
		<div
			v-if="status === 'pending'"
			class="flex h-full w-full items-center justify-center text-white"
		>
			<p>Loading paste...</p>
		</div>

		<div
			v-else-if="error"
			class="flex h-full w-full items-center justify-center text-white"
		>
			<p>Error loading paste: {{ error.message }}</p>
		</div>

		<HastebinPage
			v-else
			:item="data"
		/>
	</HastebinLayout>
</template>
