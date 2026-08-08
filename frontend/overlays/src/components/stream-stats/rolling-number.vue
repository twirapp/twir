<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
	value: string
}>()

const chars = computed(() => props.value.split(''))

function isDigit(char: string): boolean {
	return char >= '0' && char <= '9'
}

// rolling counters sweep right-to-left, like an odometer
function charDelay(index: number): string {
	return `${(chars.value.length - 1 - index) * 30}ms`
}
</script>

<template>
	<span class="rolling-number">
		<template v-for="(char, index) in chars" :key="index">
			<span v-if="isDigit(char)" class="digit-window">
				<span
					class="digit-strip"
					:style="{
						transform: `translateY(-${char}em)`,
						transitionDelay: charDelay(index),
					}"
				>
					<span v-for="n in 10" :key="n" class="digit">{{ n - 1 }}</span>
				</span>
			</span>
			<span v-else class="char">{{ char }}</span>
		</template>
	</span>
</template>

<style scoped>
.rolling-number {
	display: inline-flex;
	font-variant-numeric: tabular-nums;
}

.char {
	display: inline-block;
	height: 1em;
	line-height: 1;
}

.digit-window {
	display: inline-block;
	height: 1em;
	overflow: hidden;
	line-height: 1;
}

.digit-strip {
	display: flex;
	flex-direction: column;
	transition: transform 0.6s cubic-bezier(0.22, 1, 0.36, 1);
}

.digit {
	height: 1em;
	line-height: 1;
}

@media (prefers-reduced-motion: reduce) {
	.digit-strip {
		transition: none;
	}
}
</style>
