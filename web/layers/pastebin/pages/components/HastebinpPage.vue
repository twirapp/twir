<script setup lang="ts">
import HastebinEditor from './HastebinEditor.vue'
import HastebinToolbar from './HastebinToolbar.vue'
import HastebinViewer from './HastebinViewer.vue'

import type { PasteBinOutputDto } from '@twir/api/openapi'

const props = defineProps<{
	item?: PasteBinOutputDto
}>()

const router = useRouter()
const api = useOapi()

const content = ref<string>(props.item?.content || '')
const isEditMode = ref(!props.item)
const editorRef = ref<InstanceType<typeof HastebinEditor>>()

async function create() {
	if (!content.value?.length) return

	const req = await api.v1.pastebinCreate({
		content: content.value,
	})
	if (req.error) {
		throw req.error
	}

	await router.push(`/h/${req.data?.id}`)
}

async function duplicate() {
	if (!props.item) return

	content.value = props.item.content
	isEditMode.value = true

	// Focus the editor after it's rendered
	nextTick(() => {
		editorRef.value?.focus()
	})
}

function newPaste() {
	content.value = ''
	isEditMode.value = true

	// Focus the editor after it's rendered
	nextTick(() => {
		editorRef.value?.focus()
	})
}

// Watch for changes to isEditMode and focus editor when it becomes true
watch(isEditMode, (newValue) => {
	if (newValue) {
		nextTick(() => {
			editorRef.value?.focus()
		})
	}
})

// Computed properties for toolbar state
const canSave = computed(() => {
	return isEditMode.value && content.value?.length > 0
})

const canCopy = computed(() => {
	return !!props.item && !isEditMode.value
})

const canViewRaw = computed(() => {
	return !!props.item && !isEditMode.value
})
</script>

<template>
	<div class="h-full w-full p-4 relative">
		<HastebinToolbar
			:can-save="canSave"
			:can-copy="canCopy"
			:can-view-raw="canViewRaw"
			:item-id="props.item?.id"
			@save="create"
			@new="newPaste"
			@copy="duplicate"
		/>

		<HastebinViewer
			v-if="props.item && !isEditMode"
			:code="props.item.content"
		/>

		<HastebinEditor
			v-else
			ref="editorRef"
			v-model="content"
		/>
	</div>
</template>

<style scoped>
/* No styles needed here as they've been moved to their respective components */
</style>
