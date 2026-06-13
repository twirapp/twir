<script lang="ts" setup>

import Base from './base.vue'

defineProps<{
	createdAt: string
}>()
</script>

<template>
	<Base
		:icon="'tabler:icon-eraser'"
		:icon-color="['#ff4f4d', '#ffaaa8']"
		:created-at="createdAt"
	>
		<template #leftContent>
			Chat cleared
		</template>
	</Base>
</template>
