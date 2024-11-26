<script lang="ts" setup>
import { IconEraser } from '@tabler/icons-vue'

import Base from './base.vue'

defineProps<{
	createdAt: string
}>()
</script>

<template>
	<Base
		:icon="IconEraser"
		:icon-color="['#ff4f4d', '#ffaaa8']"
		:created-at="createdAt"
	>
		<template #leftContent>
			Chat cleared
		</template>
	</Base>
</template>
