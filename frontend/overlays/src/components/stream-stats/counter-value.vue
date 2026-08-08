<script setup lang="ts">
defineProps<{
	value: string
	animated?: boolean
}>()
</script>

<template>
	<span class="counter-value">
		<Transition v-if="animated" name="counter-value-swap" mode="out-in">
			<span :key="value">{{ value }}</span>
		</Transition>
		<template v-else>{{ value }}</template>
	</span>
</template>

<style scoped>
.counter-value {
	font-variant-numeric: tabular-nums;
}

.counter-value-swap-enter-active,
.counter-value-swap-leave-active {
	transition: opacity 0.18s ease;
}

.counter-value-swap-enter-from,
.counter-value-swap-leave-to {
	opacity: 0;
}
</style>
