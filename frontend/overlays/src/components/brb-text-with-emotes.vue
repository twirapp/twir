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
			<span
				v-else-if="chunk.type === 'emote'"
				class="emote-wrapper"
			>
				<img
					:src="chunk.value"
					class="emote"
				>
				<img
					v-for="(modifier, idx) of chunk.zeroWidthModifiers"
					:key="idx"
					:src="modifier"
					class="emote-zerowidth"
				>
			</span>
		</template>
	</span>
</template>

<style scoped>
.text-with-emotes {
	display: inline;
}

.emote-wrapper {
	position: relative;
	display: inline-block;
	margin-left: 0.25em;
	margin-right: 0.25em;
	vertical-align: middle;
}

.emote {
	display: inline-block;
	vertical-align: middle;
	height: 1.5em;
	width: auto;
	object-fit: contain;
}

.emote-zerowidth {
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	height: 1.5em;
	width: auto;
	object-fit: contain;
}
</style>
