<script setup lang="ts">
import { LoaderCircleIcon, UploadIcon } from 'lucide-vue-next'
import { useId } from 'radix-vue'
import { computed, ref } from 'vue'

import Button from './Button.vue'

defineOptions({
	inheritAttrs: false,
})
const props = withDefaults(defineProps<{
	accept?: string[] | string
	loading?: boolean
}>(), { loading: false })
const emit = defineEmits<{
	fileSelected: [File]
}>()
const isDragging = ref(false)
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

function handleFileChange(event: Event) {
	const target = event.target as HTMLInputElement
	if (!target) return

	const file = target.files?.[0]
	if (file) {
		selectedFile.value = file
		emit('fileSelected', file)
	}
}

function handleDrop(event: DragEvent) {
	isDragging.value = false
	const file = event.dataTransfer?.files[0]
	if (file) {
		selectedFile.value = file
		emit('fileSelected', file)
	}
}

function reset() {
	selectedFile.value = null
	if (fileInput.value) {
		fileInput.value.value = ''
	}
}

defineExpose({ reset })
const id = useId()

const computedAccept = computed(() => {
	if (!props.accept) return undefined

	if (Array.isArray(props.accept)) {
		return props.accept.join(',')
	}

	return props.accept
})
</script>

<template>
	<div
		class="relative w-full max-w-md"
		@dragover.prevent="isDragging = true"
		@dragleave.prevent="isDragging = false"
		@drop.prevent="handleDrop"
	>
		<Button
			v-bind="$attrs"
			variant="outline"
			class="w-full justify-center border-dashed border-2"
			as-child
		>
			<label
				:for="id"
				class="cursor-pointer flex items-center gap-2"
			>
				<template v-if="!loading">
					<UploadIcon class="size-4" />
					<!-- {{ selectedFile ? selectedFile.name : 'Upload File' }} -->
					Upload File
				</template>
				<template v-else>
					<LoaderCircleIcon class="size-10 animate-spin" />
				</template>
			</label>
		</Button>
		<input
			:id
			ref="fileInput"
			type="file"
			class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
			:accept="computedAccept"
			:disabled="loading"
			@change="handleFileChange"
		/>
	</div>
</template>
