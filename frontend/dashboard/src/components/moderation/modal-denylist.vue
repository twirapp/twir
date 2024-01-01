<script setup lang="ts">
import { IconTrash } from '@tabler/icons-vue';
import { NAlert, NInput, NDivider, NButton } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers';

const { editableItem } = useEditableItem();

const { t } = useI18n();

const addItem = () => editableItem.value!.data!.denyList.push('');
const removeItem = (i: number) => editableItem.value!.data!.denyList = editableItem.value!.data!.denyList.filter((_, idx) => idx != i);
</script>

<template>
	<div>
		<n-divider style="margin: 0; padding: 0" />

		<n-alert v-if="!editableItem?.data?.denyList.length" type="warning">
			{{ t('moderation.types.deny_list.empty') }}
		</n-alert>

		<n-alert v-else type="info" v-html="t('moderation.types.deny_list.regexp')"></n-alert>

		<div
			v-for="(_, i) of editableItem!.data!.denyList"
			:key="i"
			style="display: flex; gap: 4px;"
		>
			<n-input
				v-model:value="editableItem!.data!.denyList[i]"
				type="textarea"
				autosize
				:maxlength="500"
			/>
			<n-button text type="error" @click="removeItem(i)">
				<IconTrash />
			</n-button>
		</div>

		<n-button
			:disabled="editableItem!.data!.denyList.length >= 100"
			type="success"
			dashed
			block
			@click="addItem"
		>
			{{ t('sharedButtons.create') }}
		</n-button>
	</div>
</template>
