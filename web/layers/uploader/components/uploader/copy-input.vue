<script setup lang="ts">
const { t } = useI18n()

const props = defineProps<{
	text: string
}>()

const clipboard = useClipboard()

const copied = ref(false)
let copiedTimer: ReturnType<typeof setTimeout> | undefined

function copyText() {
	clipboard.copy(props.text)
	copied.value = true
	clearTimeout(copiedTimer)
	copiedTimer = setTimeout(() => {
		copied.value = false
	}, 2000)
}

onUnmounted(() => {
	clearTimeout(copiedTimer)
})

function selectText(event: FocusEvent) {
	;(event.target as HTMLInputElement).select()
}
</script>

<template>
	<div
		class="flex items-center rounded-xl p-1.5 border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,15%)] w-full"
	>
		<input
			:value="text"
			type="text"
			readonly
			class="flex-1 min-w-0 ml-2 mr-1 bg-transparent font-mono text-sm text-[hsl(240,11%,90%)] focus-visible:outline-none"
			@focus="selectText"
		/>
		<button
			type="button"
			class="flex-none flex items-center justify-center py-1.5 px-3 rounded-lg font-semibold border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,25%)] hover:bg-[hsl(240,11%,35%)] text-[hsl(240,11%,90%)] transition-colors"
			:title="t('uploader.actions.copyLink')"
			@click="copyText"
		>
			<Icon :name="copied ? 'lucide:check' : 'lucide:copy'" class="w-4 h-4" />
		</button>
	</div>
</template>
