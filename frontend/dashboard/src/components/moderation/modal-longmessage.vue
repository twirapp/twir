<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NInputNumber, NFormItem } from 'naive-ui';
import { useI18n } from 'vue-i18n';

defineProps<{
	item: ItemWithId
}>();

const { t } = useI18n();
</script>

<template>
	<div>
		<n-form-item :label="t('moderation.types.long_message.triggerLength')">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>
	</div>
</template>
