<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NInputNumber, NFormItem } from 'naive-ui';

defineProps<{
	item: ItemWithId
}>();
</script>

<template>
	<div>
		<n-form-item label="trigger length" feedback="qwe">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-form-item label="Max" feedback="qwe">
			<n-input-number
				v-model:value="item.data!.maxPercentage"
				:max="500"
			/>
		</n-form-item>
	</div>
</template>
