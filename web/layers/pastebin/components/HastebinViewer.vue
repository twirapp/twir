<script setup lang="ts">
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const { currentPaste } = storeToRefs(usePasteStore())
const { detectLanguage, highlight } = useHighlight()

const code = computed(() => {
	if (!currentPaste.value) return null

	const lang = detectLanguage(currentPaste.value.content)
	return highlight(currentPaste.value.content, lang)
})
</script>

<template>
	<pre class="h-full"><code v-html="code"></code></pre>
</template>

<style scoped>
@reference "~/assets/css/tailwind.css";

:deep(code) {
	font-family: 'JetBrains Mono', serif;
}
</style>
