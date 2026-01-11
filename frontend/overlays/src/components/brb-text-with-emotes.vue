<script setup lang="ts">
import type { TextChunk } from '@/composables/brb/use-brb-text-parser.js'

interface Props {
	chunks: TextChunk[]
}

defineProps<Props>()
</script>

<template>
	<span class="text-with-emotes">
		<template v-for="(chunk, index) in chunks" :key="index">
			<span v-if="chunk.type === 'text'">{{ chunk.value }}</span>
			<img
				v-else-if="chunk.type === 'emote'"
				:src="chunk.value"
				class="emote"
			>
		</template>
	</span>
</template>

<style scoped>
.text-with-emotes {
	display: inline-flex;
	align-items: center;
	gap: 0.25em;
	flex-wrap: wrap;
}

.emote {
	display: inline-block;
	vertical-align: middle;
	height: 1.5em;
	width: auto;
	object-fit: contain;
}
</style>
