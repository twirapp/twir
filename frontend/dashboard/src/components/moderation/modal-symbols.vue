<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NFormItem, NInputNumber, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

defineProps<{
	item: ItemWithId
}>();

const { t } = useI18n();
</script>

<template>
	<div>
		<n-form-item :label="t('moderation.types.symbols.triggerLength')">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-form-item :label="t('moderation.types.symbols.maxPercentage')">
			<n-input-number
				v-model:value="item.data!.maxPercentage"
				:max="100"
			/>
		</n-form-item>

		<n-divider style="margin: 0; padding: 0" />
	</div>
</template>
