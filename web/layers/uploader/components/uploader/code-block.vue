<script setup lang="ts">
const props = defineProps<{
	code: string
	/** What actually lands in the clipboard when it differs from the visible code (e.g. unmasked secrets). */
	copyText?: string
}>()

const clipboard = useClipboard()

const copied = ref(false)
let copiedTimer: ReturnType<typeof setTimeout> | undefined

function copyCode() {
	clipboard.copy(props.copyText ?? props.code)
	copied.value = true
	clearTimeout(copiedTimer)
	copiedTimer = setTimeout(() => {
		copied.value = false
	}, 2000)
}

onUnmounted(() => {
	clearTimeout(copiedTimer)
})
</script>

<template>
	<div class="relative rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)]">
		<pre class="overflow-x-auto p-3 text-xs leading-relaxed text-[hsl(240,11%,80%)]"><code>{{ code }}</code></pre>
		<button
			type="button"
			class="absolute top-2 right-2 flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
			:title="$t('uploader.guide.copy')"
			@click="copyCode"
		>
			<Icon :name="copied ? 'lucide:check' : 'lucide:copy'" class="h-3.5 w-3.5" />
		</button>
	</div>
</template>
