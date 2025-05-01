<script setup lang="ts">
import { usePasteStore } from '../../stores/pasteStore'

const textareaRef = ref<HTMLTextAreaElement>()
const pasteStore = usePasteStore()

// Focus textarea when component is mounted
onMounted(() => {
	textareaRef.value?.focus()
})

// Expose focus method to parent component
defineExpose({
	focus: () => {
		nextTick(() => {
			textareaRef.value?.focus()
		})
	},
})
</script>

<template>
	<textarea
		ref="textareaRef"
		v-model="pasteStore.editableContent"
		class="h-full w-full p-2 bg-transparent outline-none rounded-md input"
	/>
</template>

<style scoped>
.input {
  font-family: 'JetBrains Mono';
}
</style>
