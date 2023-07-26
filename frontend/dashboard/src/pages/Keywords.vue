<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import type { Keyword } from '@twir/grpc/generated/api/api/keywords';
import { useThrottleFn } from '@vueuse/core';
import {
	type DataTableColumns,
  NDataTable,
  NSpace,
  NTag,
  NSwitch,
  NButton,
  NPopconfirm,
  NModal,
	NResult,
} from 'naive-ui';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useKeywordsManager, useUserAccessFlagChecker } from '@/api/index.js';
import Modal from '@/components/keywords/modal.vue';
import { type EditableKeyword } from '@/components/keywords/types.js';
import { renderIcon } from '@/helpers/index.js';

const keywordsManager = useKeywordsManager();
const keywords = keywordsManager.getAll({});
const keywordsDeleter = keywordsManager.deleteOne;
const keywordsPatcher = keywordsManager.patch!;

const throttledSwitchState = useThrottleFn((id: string, v: boolean) => {
	keywordsPatcher.mutate({ id, enabled: v });
}, 500);

const showModal = ref(false);

const userCanViewKeywords = useUserAccessFlagChecker('VIEW_KEYWORDS');
const userCanManageKeywords = useUserAccessFlagChecker('MANAGE_KEYWORDS');

const { t } = useI18n();

const columns = computed<DataTableColumns<Keyword>>(() => [
	{
		title: t('keywords.triggerText'),
		key: 'text',
		render(row) {
			return h(
				NTag,
				{ type: 'info', bordered: false },
				{
					default: () => row.text.slice(0, 100) + (row.text.length > 100 ? '...' : ''),
				},
				);
		},
	},
	{
		title: t('sharedTexts.response'),
		key: 'response',
		width: 200,
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.response || 'No response' });
		},
	},
	{
		title: t('keywords.usages'),
		key: 'usages',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.usages });
		},
	},
	{
		title: t('sharedTexts.status'),
		key: 'enabled',
		render(row) {
			return h(NSwitch, {
				value: row.enabled,
				onUpdateValue: (value) => {
					throttledSwitchState(row.id, value);
				},
				disabled: !userCanManageKeywords.value,
			});
		},
	},
	{
		title: t('sharedTexts.actions'),
		key: 'actions',
		width: 150,
		render(row) {
			const editButton = h(
				NButton,
				{
					type: 'primary',
					size: 'small',
					onClick: () => openModal(row),
					quaternary: true,
					disabled: !userCanManageKeywords.value,
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => keywordsDeleter.mutate({ id: row.id }),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
						disabled: !userCanManageKeywords.value,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => t('deleteConfirmation.text'),
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
]);

const editableKeyword = ref<EditableKeyword | null>(null);
function openModal(t: EditableKeyword | null) {
	editableKeyword.value = t;
	showModal.value = true;
}
function closeModal() {
	showModal.value = false;
}
</script>

<template>
	<n-result v-if="!userCanViewKeywords" status="403" :title="t('haveNoAccess.message')" />

	<div v-else>
		<n-space justify="space-between" align="center">
			<h2>{{ t('keywords.title') }}</h2>
			<n-button :disabled="!userCanManageKeywords" secondary type="success" @click="openModal(null)">
				{{ t('sharedButtons.create') }}
			</n-button>
		</n-space>
		<n-data-table
			:isLoading="keywords.isLoading.value"
			:columns="columns"
			:data="keywords.data.value?.keywords ?? []"
		/>

		<n-modal
			v-model:show="showModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="editableKeyword?.id ? 'Edit keyword' : 'New keyword'"
			class="modal"
			:style="{
				width: '600px',
				top: '50px',
			}"
			:on-close="closeModal"
		>
			<modal :keyword="editableKeyword" @close="closeModal" />
		</n-modal>
	</div>
</template>
