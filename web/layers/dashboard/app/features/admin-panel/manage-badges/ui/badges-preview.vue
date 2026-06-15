<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ image: File | string }>()
const imageSrc = computed(() => {
	if (props.image instanceof File) {
		return URL.createObjectURL(props.image)
	}

	return props.image
})

const BADGE_SIZES = [144, 72, 36, 18]
</script>

<template>
	<div class="flex items-end gap-2 w-full">
		<div v-for="size of BADGE_SIZES" :key="size" class="bg-muted p-2 rounded-md">
			<img :src="imageSrc" :width="size" :height="size" />
		</div>
	</div>
</template>
