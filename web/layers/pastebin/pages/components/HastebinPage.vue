<script setup lang="ts">
import HastebinEditor from './HastebinEditor.vue'
import HastebinToolbar from './HastebinToolbar.vue'
import HastebinViewer from './HastebinViewer.vue'
import { usePasteStore } from '../../stores/pasteStore'

const editorRef = ref<InstanceType<typeof HastebinEditor>>()
const router = useRouter()
const api = useOapi()
const pasteStore = usePasteStore()

async function create() {
	if (!pasteStore.editableContent) return

	const req = await api.v1.pastebinCreate({
		content: pasteStore.editableContent,
	})
	if (req.error) {
		throw req.error
	}

	await router.push(`/h/${req.data?.id}`)
}

async function duplicate() {
	if (!pasteStore.currentPaste) return

	pasteStore.duplicate()

	console.log('here')
	// Redirect to the new paste page
	await router.push('/h')
}

function newPaste() {
	pasteStore.editableContent = ''

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
			v-if="pasteStore.currentPaste"
		/>

		<HastebinEditor
			v-else
			ref="editorRef"
		/>
	</div>
</template>
