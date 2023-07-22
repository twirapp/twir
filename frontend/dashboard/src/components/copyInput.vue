<script setup lang='ts'>
import { useTimeout } from '@vueuse/core';
import {
  NInputGroup,
  NInput,
  NButton,
} from 'naive-ui';

import { copyToClipBoard } from '@/helpers/index.js';

const props = defineProps<{
	text: string,
}>();

const { start: copyStart, isPending } = useTimeout(2000, { controls: true, immediate: false });

async function copy() {
	await copyToClipBoard(props.text);
	copyStart();
}
</script>

<template>
	<n-input-group>
		<n-input type="text" size="small" :value="text" disabled />
		<n-button size="small" type="primary" @click="copy">
			<span v-if="!isPending">Copy</span>
			<span v-else>Copied</span>
		</n-button>
	</n-input-group>
</template>
