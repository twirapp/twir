<script setup lang="ts">
import { IconTrash } from '@tabler/icons-vue';
import { NAlert, NInput, NDivider, NButton, NA } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers';

const { editableItem } = useEditableItem();

const { t } = useI18n();

const addItem = () => editableItem.value!.data!.denyList.push('');
const removeItem = (i: number) => editableItem.value!.data!.denyList = editableItem.value!.data!.denyList.filter((_, idx) => idx != i);
</script>

<template>
	<div>
		<n-divider class="m-0 p-0" />

		<n-alert v-if="!editableItem?.data?.denyList.length" type="warning">
			{{ t('moderation.types.deny_list.empty') }}
		</n-alert>

		<n-alert v-else type="info">
			<i18n-t
				keypath="moderation.types.deny_list.regexp"
			>
				<n-a
					href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
					target="_blank"
				>
					{{ t('moderation.types.deny_list.regexpCheatSheet') }}
				</n-a>
			</i18n-t>
		</n-alert>

		<div
			v-for="(_, i) of editableItem!.data!.denyList"
			:key="i"
			class="flex gap-1"
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
