<script setup lang="ts">
import HastebinPage from '#layers/pastebin/pages/components/hastebin-page.vue'
import { useOapi } from '~/composables/use-oapi'

const route = useRoute()
const api = useOapi()

const { data } = await useAsyncData(
	'hastebin',
	async () => {
		const req = await api.v1.pastebinGetById(route.params.id as string)
		if (req.error) {
			throw req.error
		}

		return req.data
	},
)
</script>

<template>
	<div class="flex h-screen bg-[#1E1E1E]">
		<HastebinPage :item="data" />
	</div>
</template>
