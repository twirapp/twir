<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NInputNumber, NFormItem } from 'naive-ui';

defineProps<{
	item: ItemWithId
}>();
</script>

<template>
	<div>
		<n-form-item label="Max emotes" feedback="qwe">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>
	</div>
</template>
