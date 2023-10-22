<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NFormItem, NInputNumber, NDivider } from 'naive-ui';

defineProps<{
	item: ItemWithId
}>();
</script>

<template>
	<div>
		<n-form-item label="Trigger length" feedback="from what message length start to check on symbols">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-form-item label="Max percentage" feedback="qwe">
			<n-input-number
				v-model:value="item.data!.maxPercentage"
				:max="500"
			/>
		</n-form-item>

		<n-divider style="margin: 0; padding: 0" />
	</div>
</template>
