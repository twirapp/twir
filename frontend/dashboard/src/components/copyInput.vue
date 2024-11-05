<script setup lang='ts'>
import { useTimeout } from '@vueuse/core'
import {
	NButton,
	NInput,
	NInputGroup,
} from 'naive-ui'

import type { Size } from 'naive-ui/es/input/src/interface.js'

import { copyToClipBoard } from '@/helpers/index.js'

const props = withDefaults(defineProps<{
	text: string
	type?: 'password' | 'text'
	size?: Size
}>(), {
	type: 'text',
	size: 'small',
})

const { start: copyStart, isPending } = useTimeout(2000, { controls: true, immediate: false })

async function copy() {
	await copyToClipBoard(props.text)
	copyStart()
}
</script>

<template>
	<NInputGroup>
		<NInput :type="type" show-password-on="click" size="small" :value="text" @update-value="() => {}" />
		<NButton :size="size" type="primary" @click="copy">
			<span v-if="!isPending">Copy</span>
			<span v-else>Copied</span>
		</NButton>
	</NInputGroup>
</template>
