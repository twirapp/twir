<script setup lang="ts">
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const { currentPaste } = storeToRefs(usePasteStore())
</script>

<template>
	<Shiki
		v-if="currentPaste?.content"
		:code="currentPaste.content"
		class="h-full"
	/>
</template>

<style scoped>
:deep(pre code) {
  font-family: 'JetBrains Mono';
}

:deep(code) {
	@apply break-words text-wrap
}
</style>
