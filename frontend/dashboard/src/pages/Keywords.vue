<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import {
	type DataTableColumns,
	NDataTable,
	NSpace,
	NTag,
	NSwitch,
	NButton,
	NPopconfirm,
	NModal,
} from 'naive-ui';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';
import { useKeywordsApi, type Keyword } from '@/api/keywords';
import Modal from '@/components/keywords/modal.vue';
import { renderIcon } from '@/helpers/index.js';

const { t } = useI18n();
const userCanManageKeywords = useUserAccessFlagChecker('MANAGE_KEYWORDS');
const showModal = ref(false);

const keywordsApi = useKeywordsApi();
const keywordsUpdate = keywordsApi.useMutationUpdateKeyword();
const keywordsRemove = keywordsApi.useMutationRemoveKeyword();

const { data, fetching } = keywordsApi.useQueryKeywords();
const keywords = computed(() => {
	return data.value?.keywords ?? [];
});

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
			return h(NTag, {
				type: 'info',
				bordered: true,
			}, { default: () => row.response || 'No response' });
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
				onUpdateValue: (enabled) => {
					keywordsUpdate.executeMutation({ id: row.id, opts: { enabled } });
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
					onPositiveClick: () => keywordsRemove.executeMutation({ id: row.id }),
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

			return h(NSpace, {}, { default: () => [editButton, deleteButton] });
		},
	},
]);

const editableKeyword = ref<Keyword | null>(null);

function openModal(keyword: Keyword | null) {
	editableKeyword.value = keyword;
	showModal.value = true;
}

function closeModal() {
	showModal.value = false;
}
</script>

<template>
	<n-space justify="space-between" align="center">
		<h2>{{ t('keywords.title') }}</h2>
		<n-button :disabled="!userCanManageKeywords" secondary type="success" @click="openModal(null)">
			{{ t('sharedButtons.create') }}
		</n-button>
	</n-space>

	<n-data-table
		:isLoading="fetching"
		:columns="columns"
		:data="keywords"
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
</template>
