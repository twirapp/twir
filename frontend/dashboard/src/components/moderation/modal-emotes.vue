<script setup lang="ts">
import { NInputNumber, NFormItem, NDivider } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers';

const { editableItem } = useEditableItem();
const { t } = useI18n();
</script>

<template>
	<div>
		<n-form-item :label="t('moderation.types.emotes.triggerLength')">
			<n-input-number
				v-model:value="editableItem!.data!.triggerLength"
				:max="500"
			/>
		</n-form-item>

		<n-divider class="m-0" />
	</div>
</template>
