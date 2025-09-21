<script lang="ts" setup>
import { Button } from '@/components/ui/button'
import { useEyeDropper } from '@vueuse/core'
import { BanIcon, PipetteIcon } from 'lucide-vue-next'
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
	<Button
		v-if="isSupported"
		@click="handleEyeDropperClick"
		type="button"
		:class="props.class"
		aria-label="Pick color from screen"
		title="Pick color from screen"
		variant="outline"
		size="none"
	>
		<PipetteIcon class="w-4 h-4" />
	</Button>

	<Button
		v-else
		:disabled="true"
		variant="outline"
		size="none"
		:class="props.class"
		title="EyeDropper API не поддерживается в вашем браузере"
	>
		<BanIcon class="w-4 h-4" />
	</Button>
</template>
