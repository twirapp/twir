<script setup lang="ts">
import { IconMenu2, IconX } from '@tabler/icons-vue';

defineProps<{
	isOpen: boolean
}>();
</script>

<template>
	<div class="hamburger">
		<icon-menu2 v-if="!isOpen" />
		<icon-x v-else />
	</div>
</template>

<style scoped>
.hamburger {
	cursor: pointer;
	width: 38px;
	height: 38px !important;
	display: flex;
	align-items: center !important;
	justify-content: center !important;
}
</style>
