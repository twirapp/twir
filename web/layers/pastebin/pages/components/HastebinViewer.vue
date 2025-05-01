<script setup lang="ts">
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const { currentPaste } = storeToRefs(usePasteStore())
</script>

<template>
	<Shiki
		v-if="currentPaste?.content"
		:code="currentPaste.content"
		as="div"
		class="h-full"
	/>
</template>

<style scoped>
:deep(.shiki) {
  @apply h-full w-full;
}

:deep(code) {
  counter-reset: step;
  counter-increment: step 0;
}

:deep(pre code) {
  font-family: 'JetBrains Mono';
}

:deep(code .line::before) {
  content: counter(step);
  counter-increment: step;
  width: 1rem;
  margin-right: 1.5rem;
  display: inline-block;
  text-align: right;
  color: rgba(115, 138, 148, .4);
}
</style>
