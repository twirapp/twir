<script setup>
import Delimiter from '@editorjs/delimiter'
import EditorJS from '@editorjs/editorjs'
import Header from '@editorjs/header'
import List from '@editorjs/list'
import Paragraph from '@editorjs/paragraph'
import Quote from '@editorjs/quote'
import SimpleImage from '@editorjs/simple-image'
import Underline from '@editorjs/underline'
import { onMounted, onUnmounted, ref } from 'vue'
import './notifications-form.css'

const modelValue = defineModel()

const editor = ref(null)
const editorHtmlElement = ref(null)

onMounted(() => {
	if (!editorHtmlElement.value) return

	editor.value = new EditorJS({
		holder: editorHtmlElement.value,
		onChange() {
			editor.value?.save().then((outputData) => {
				modelValue.value = JSON.stringify(outputData)
			})
		},
		onReady() {
			if (!modelValue.value) {
				return
			}

			const data = JSON.parse(modelValue.value)
			editor.value?.blocks.render(data)
		},
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
	})
})

onUnmounted(() => {
	editor.value.destroy()
})
</script>

<template>
	<div ref="editorHtmlElement" class="bg-background border border-border rounded-md"></div>
</template>
