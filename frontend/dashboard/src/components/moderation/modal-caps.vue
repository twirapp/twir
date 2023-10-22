<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NInputNumber, NFormItem, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

defineProps<{
	item: ItemWithId
}>();

const { t } = useI18n();
</script>

<template>
	<div>
		<n-form-item :label="t('moderation.types.caps.triggerLength')">
			<n-input-number
				v-model:value="item.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-form-item :label="t('moderation.types.caps.maxPercentage')">
			<n-input-number
				v-model:value="item.data!.maxPercentage"
				:max="100"
			/>
		</n-form-item>

		<n-divider style="margin: 0px; padding: 0px;" />
	</div>
</template>
