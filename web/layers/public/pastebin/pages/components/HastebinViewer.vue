<script setup lang="ts">
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const { currentPaste } = storeToRefs(usePasteStore())
</script>

<template>
	<Shiki v-if="currentPaste?.content" :code="currentPaste.content" class="h-full" />
</template>

<style scoped>
@reference '~/assets/css/tailwind.css';

:deep(pre code) {
	font-family: 'JetBrains Mono', serif;
}

:deep(code) {
	@apply wrap-break-word text-wrap;
}
</style>
