<script setup lang="ts">
import { NInputNumber, NFormItem, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers';

const { editableItem } = useEditableItem();

const { t } = useI18n();
</script>

<template>
	<div>
		<n-form-item :label="t('moderation.types.caps.triggerLength')">
			<n-input-number
				v-model:value="editableItem!.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-form-item :label="t('moderation.types.caps.maxPercentage')">
			<n-input-number
				v-model:value="editableItem!.data!.maxPercentage"
				:max="100"
			/>
		</n-form-item>

		<n-divider class="m-0 p-0" />
	</div>
</template>
