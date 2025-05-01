<script setup lang="ts">
import HastebinEditor from './HastebinEditor.vue'
import HastebinToolbar from './HastebinToolbar.vue'
import HastebinViewer from './HastebinViewer.vue'
import { usePasteStore } from '../../stores/pasteStore'

const editorRef = ref<InstanceType<typeof HastebinEditor>>()
const router = useRouter()
const api = useOapi()
const store = usePasteStore()
const { currentPaste, editableContent } = storeToRefs(store)

async function create() {
	if (!editableContent.value) return

	const req = await api.v1.pastebinCreate({
		content: editableContent.value,
	})
	if (req.error) {
		throw req.error
	}

	await router.push(`/h/${req.data?.id}`)
}

async function duplicate() {
	if (!currentPaste.value) return

	store.duplicate()
	await router.push('/h')
}

function newPaste() {
	editableContent.value = ''

	// Focus the editor after it's rendered
	nextTick(() => {
		editorRef.value?.focus()
	})
}
</script>

<template>
	<div class="h-full w-full p-4 relative">
		<HastebinToolbar
			@save="create"
			@new="newPaste"
			@copy="duplicate"
		/>

		<HastebinViewer
			v-if="currentPaste"
		/>

		<HastebinEditor
			v-else
			ref="editorRef"
		/>
	</div>
</template>
