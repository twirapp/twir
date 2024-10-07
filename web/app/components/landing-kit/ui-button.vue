<script setup lang="ts">
import { computed } from 'vue'

interface Props {
	variant: 'primary' | 'secondary'
	href: string
}

const props = defineProps<Props>()

const variantStyles: Record<Props['variant'], string> = {
	primary: 'bg-[#5D58F5] hover:bg-[#6964FF] focus-visible:ring-[#5D58F5]/50',
	secondary: 'ring-[#393A3E] ring-1 hover:bg-[#393A3E]/50',
}

const variant = computed(() => variantStyles[props.variant])
</script>

<template>
	<a
		:href="href" class="inline-flex justify-center items-center py-3 px-7 xs:py-4 font-semibold text-white rounded-lg transition-[background,box-shadow] text-base sm:text-lg focus-visible:outline-none focus-visible:ring-4 text-center whitespace-nowrap" :class="[
			variant,
		]"
	>
		<slot></slot>
	</a>
</template>
