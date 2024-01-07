<script setup lang="ts">
import { NButton, useThemeVars } from 'naive-ui';
import { computed } from 'vue';

const props = defineProps<{
	url: string
}>();

const themeVars = useThemeVars();
const iconColor = computed(() => themeVars.value.textColor2);
const openLink = () => window.open(props.url, '_blank');
</script>

<template>
	<n-button
		quaternary
		circle
		style="padding: 5px"
		@click="openLink"
	>
		<slot></slot>
	</n-button>
</template>

<style scoped>
:slotted(svg) {
  width: 20px;
	color: v-bind('iconColor');
	fill: currentColor;
}
</style>
