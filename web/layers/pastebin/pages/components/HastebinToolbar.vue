<script setup lang="ts">
import { Icon } from '#components'
import { usePasteStore } from '#layers/pastebin/stores/pasteStore'

const emit = defineEmits<{
	save: []
	new: []
	copy: []
}>()

const pasteStore = usePasteStore()
const { currentPaste, editableContent } = storeToRefs(pasteStore)

const requestUrl = useRequestURL()

const canSave = computed(() => {
	if (currentPaste.value) {
		return false
	}

	return editableContent.value?.length
})

const canCopy = computed(() => {
	return !!currentPaste.value
})

const buttons = computed(() => {
	return [
		{
			name: 'save',
			icon: 'lucide:save',
			disabled: !canSave.value,
			tooltip: 'Save',
			onClick: () => emit('save'),
		},
		{
			name: 'new',
			icon: 'lucide:file-plus',
			tooltip: 'New',
			onClick: () => emit('new'),
		},
		{
			name: 'copy',
			icon: 'lucide:copy',
			disabled: !canCopy.value,
			tooltip: 'Copy',
			onClick: () => emit('copy'),
		},
		{
			name: 'text',
			icon: 'lucide:text',
			disabled: !currentPaste.value,
			tooltip: 'Text',
			href: currentPaste.value?.id ? `${requestUrl.origin}/h/${currentPaste.value?.id}/raw` : undefined,
		},
	]
})
</script>

<template>
	<div class="flex flex-row gap-2 items-center absolute top-2 right-4 h-10 p-2 bg-gray-500/50 rounded-md">
		<UiButton
			v-for="button of buttons"
			:key="button.name"
			variant="link"
			size="icon"
			:disabled="button.disabled"
			:as="button.href ? 'a' : 'button'"
			:href="button.href"
			@click="button.onClick"
		>
			<Icon class="icon" :name="button.icon" />
		</UiButton>
	</div>
</template>

<style scoped>
.icon {
  @apply size-6 cursor-pointer text-gray-300;
}

.icon:disabled {
  @apply opacity-50;
}

.icon:hover {
  @apply text-gray-100;
}
</style>
