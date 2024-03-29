<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { type Timer } from '@twir/api/messages/timers/timers';
import {
	type DataTableColumns,
	NButton,
	NDataTable,
	NModal,
	NPopconfirm,
	NSpace,
	NSwitch,
	NTag,
} from 'naive-ui';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useTimersManager, useUserAccessFlagChecker } from '@/api/index.js';
import ImportModal from '@/components/timers/importModal.vue';
import Modal from '@/components/timers/modal.vue';
import { type EditableTimer } from '@/components/timers/types.js';
import { renderIcon } from '@/helpers/index.js';

const timersManager = useTimersManager();
const timers = timersManager.getAll({});
const timersDeleter = timersManager.deleteOne;
const timersPatcher = timersManager.patch!;

const userCanManageTimers = useUserAccessFlagChecker('MANAGE_TIMERS');

const { t } = useI18n();

const columns = computed<DataTableColumns<Timer>>(() => [
	{
		title: 'Name',
		key: 'name',
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, { default: () => row.name });
		},
	},
	{
		title: t('sharedTexts.responses'),
		key: 'responses',
		render(row) {
			return h(NSpace, { vertical: true }, {
				default: () => row.responses.map((response) => {
					return h('span', null, response.text);
				}),
			});
		},
	},
	{
		title: t('timers.table.columns.intervalInMinutes'),
		key: 'timeInterval',
		render(row) {
			return h(NTag, { type: 'info' }, { default: () => `${row.timeInterval} m.` });
		},
	},
	// {
	// 	title: 'Interval in messages',
	// 	key: 'messageInterval',
	// 	render(row) {
	// 		return h(NTag, { type: 'info' }, { default: () => `${row.messageInterval}` });
	// 	},
	// },
	{
		title: t('sharedTexts.status'),
		key: 'enabled',
		render(row) {
			return h(NSwitch, {
				value: row.enabled,
				disabled: !userCanManageTimers.value,
				onUpdateValue: (v) => {
					timersPatcher.mutate({ id: row.id!, enabled: v });
				},
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
					disabled: !userCanManageTimers.value,
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => timersDeleter.mutate({ id: row.id }),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
						disabled: !userCanManageTimers.value,
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

const showModal = ref(false);

const editableTimer = ref<EditableTimer | null>(null);

const showImportModal = ref(false);

function openModal(t: EditableTimer | null) {
	editableTimer.value = t;
	showModal.value = true;
}

function closeModal() {
	showModal.value = false;
}

const timersLength = computed(() => timers.data?.value?.timers?.length ?? 0);
</script>

<template>
	<n-space justify="space-between" align="center">
		<h2>{{ t('sidebar.timers') }}</h2>
		<div>
			<n-button
				:disabled="!userCanManageTimers" secondary type="info"
				@click="showImportModal = true"
			>
				Import
			</n-button>
			<n-button
				secondary type="success"
				:disabled="!userCanManageTimers || timersLength >= 10"
				@click="openModal(null)"
			>
				{{ timersLength >= 10 ? t('timers.limitExceeded') : t('sharedButtons.create') }} ({{
					timersLength }}/10)
			</n-button>
		</div>
	</n-space>
	<n-data-table
		:isLoading="timers.isLoading.value"
		:columns="columns"
		:data="timers.data.value?.timers ?? []"
	/>

	<n-modal
		v-model:show="showImportModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Import timers"
		class="modal"
		:style="{
			width: '600px',
			top: '50px',
		}"
		:on-close="() => showImportModal = false"
	>
		<import-modal @close="showImportModal = false" />
	</n-modal>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableTimer?.name ?? 'New timer'"
		class="modal"
		:style="{
			width: '600px',
			top: '50px',
		}"
		:on-close="closeModal"
	>
		<modal :timer="editableTimer" @close="closeModal" />
	</n-modal>
</template>
