<script setup lang="ts">
import { Toggle } from 'radix-vue';

const pressed = defineModel<boolean>({ required: true });
</script>

<template>
	<Toggle
		v-model:pressed="pressed"
		style="outline: none"
	>
		<div class="inline-flex h-6 w-6 items-center justify-center">
			<div :data-state="pressed ? 'on' : 'off'" class="toggle"></div>
		</div>
	</Toggle>
</template>

<style scoped>
.toggle {
	height: 12px;
	width: 22px;
	position: relative;
}

.toggle::after,
.toggle::before {
	content: "";
	position: absolute;
	display: block;
	width: 100%;
	height: 1.75px;
	background-color: currentColor;
	border-radius: 1px;
	transition: all 0.2s cubic-bezier(0.1, 0.82, 0.76, 0.965);
}

.toggle::before {
	top: 0;
}

.toggle::after {
	bottom: 0;
}

.toggle[data-state="on"]::before {
	transform: rotate(45deg);
	top: 5px;
}

.toggle[data-state="on"]::after {
	transform: rotate(-45deg);
	bottom: 5px;
}
</style>
