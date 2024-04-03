<script setup lang="ts">
import { NTooltip, NTime } from 'naive-ui';

defineProps<{ time: Date }>();
</script>

<template>
	<n-tooltip placement="left" trigger="hover">
		<template #trigger>
			<n-time class="cursor-pointer" :time="time" type="relative" />
		</template>
		<span>{{ time.toLocaleString() }}</span>
	</n-tooltip>
</template>
