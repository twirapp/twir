<script lang="ts" setup>
import { useEyeDropper } from '@vueuse/core'
import tinycolor from 'tinycolor2'

import type { HTMLAttributes } from 'vue'

const props = defineProps<{
	class?: HTMLAttributes['class']
}>()

const emit = defineEmits<{
	'pick-color': [color: string]
}>()

const { isSupported, open } = useEyeDropper()

async function handleEyeDropperClick() {
	try {
		const result = await open()
		if (result && result.sRGBHex) {
			const color = tinycolor(result.sRGBHex)
			emit('pick-color', color.toHex())
		}
	} catch (error) {
		console.error('EyeDropper error:', error)
	}
}
</script>

<template>
	<UiButton
		v-if="isSupported"
		type="button"
		:class="props.class"
		aria-label="Pick color from screen"
		title="Pick color from screen"
		variant="outline"
		size="custom"
		@click="handleEyeDropperClick"
	>
		<Icon name="lucide:pipette" class="size-4" />
	</UiButton>

	<UiButton
		v-else
		:disabled="true"
		variant="outline"
		size="custom"
		:class="props.class"
		title="EyeDropper API don't supported in this browser"
	>
		<Icon name="lucide:ban" class="size-4" />
	</UiButton>
</template>
