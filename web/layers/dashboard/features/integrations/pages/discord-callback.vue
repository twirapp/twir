<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { LoaderCircleIcon } from 'lucide-vue-next'

import { useDiscordIntegration } from '../composables/discord/use-discord-integration.ts'

const route = useRoute()
const { connectGuild } = useDiscordIntegration()

onMounted(async () => {
	if (typeof route.query.code === 'string') {
		await connectGuild(route.query.code)
		window.opener.postMessage('discord-connected')
		window.close()
	}
})
</script>

<template>
	<div class="flex justify-center items-center h-full">
		<LoaderCircleIcon class="animate-spin size-12" />
	</div>
</template>
