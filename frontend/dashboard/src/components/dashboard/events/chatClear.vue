<script lang="ts" setup>
import { IconEraser } from '@tabler/icons-vue';

import Base from './base.vue';

defineProps<{
	createdAt: string,
}>();
</script>

<template>
	<Base :icon="IconEraser" :created-at="createdAt">
		<template #leftContent>
			<div style="display: flex; flex-direction:column">
				<span>Chat Cleared</span>
			</div>
		</template>
	</Base>
</template>
