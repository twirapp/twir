<script setup lang="ts">
import { IconMenu2, IconX } from '@tabler/icons-vue';

defineProps<{
	isOpen: boolean
}>();
</script>

<template>
	<div class="flex cursor-pointer h-9 w-9 items-center justify-center">
		<icon-menu2 v-if="!isOpen" />
		<icon-x v-else />
	</div>
</template>
