<!-- eslint-disable vue/no-v-html -->
<!-- eslint-disable no-undef -->
<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core';
import { ref, watch, nextTick, computed } from 'vue';

import { useOverlaysParseHtml } from '@/api/registry';

const props = defineProps<{
	index: number | string;
	posX: number;
	posY: number;
	width: number;
	height: number;
	text: string;
	css: string;
	js: string
	periodicallyRefetchData: boolean,
}>();

const fetcher = useOverlaysParseHtml();

const exampleValue = ref('');

const { pause, resume } = useIntervalFn(async () => {
	const data = await fetcher.mutateAsync(props.text ?? '');
	exampleValue.value = data ?? '';
}, 1000, { immediate: true, immediateCallback: true });

const executeFunc = computed(() => {
	return new Function(`${props.js}; onDataUpdate();`);
});

watch(exampleValue, async () => {
	await nextTick();
	// calling user defined function
	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-ignore
	executeFunc.value?.();
});

watch(props, (p) => {
	const v = p.periodicallyRefetchData;

	if (!v) pause();
	else resume();
}, { immediate: true });
</script>

<template>
	<div
		:id="'layer-' + index"
		style="position: absolute;"
		:style="{
			transform: `translate(${posX}px, ${posY}px)`,
			width: `${width}px`,
			height: `${height}px`,
			'text-wrap': 'nowrap',
			overflow: 'hidden'
		}"
	>
		<component :is="'style'">
			{{
				`#layersExampleRender${index} {
					${css}
				}`
			}}
		</component>

		<div :id="'layersExampleRender'+index" :style="css" v-html="exampleValue" />
	</div>
</template>
