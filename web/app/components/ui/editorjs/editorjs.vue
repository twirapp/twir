<script setup>
import Delimiter from '@editorjs/delimiter'
import EditorJS from '@editorjs/editorjs'
import Header from '@editorjs/header'
import List from '@editorjs/list'
import Paragraph from '@editorjs/paragraph'
import Quote from '@editorjs/quote'
import SimpleImage from '@editorjs/simple-image'
import Underline from '@editorjs/underline'
import { onMounted, onUnmounted, ref, watch } from 'vue'
import './notifications-form.css'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

let updatingModel = false

const editor = ref(null)
const editorHtmlElement = ref(null)

// model -> view
function modelToView() {
	if (!props.modelValue) {
		editor.value.blocks.clear()
		return
	}

	editor.value.blocks.render(JSON.parse(props.modelValue))
}
// view -> model
function viewToModel(api, event) {
	updatingModel = true
	editor.value.save().then((outputData) => {
		emit('update:modelValue', JSON.stringify(outputData))
	}).catch((error) => {
		console.log(event, 'Saving failed: ', error)
	}).finally(() => {
		updatingModel = false
	})
}

onMounted(() => {
	if (!editorHtmlElement.value) return

	editor.value = new EditorJS({
		holder: editorHtmlElement.value,
		onReady: modelToView,
		onChange: viewToModel,
		placeholder: 'Type here...',
		tools: {
			header: Header,
			image: {
				class: SimpleImage,
				inlineToolbar: true,
				config: {
					placeholder: 'Paste image URL',
				},
			},
			list: {
				class: List,
				inlineToolbar: true,
				config: {
					defaultStyle: 'unordered',
				},
			},
			delimiter: Delimiter,
			paragraph: {
				class: Paragraph,
				inlineToolbar: true,
			},
			quote: Quote,
			underline: Underline,
		},
		data: props.modelValue,
	})
})

watch(() => props.modelValue, () => {
	if (!updatingModel) {
		modelToView()
	}
})

onUnmounted(() => {
	editor.value.destroy()
})
</script>

<template>
	<div ref="editorHtmlElement" class="bg-background border border-border rounded-md"></div>
</template>
