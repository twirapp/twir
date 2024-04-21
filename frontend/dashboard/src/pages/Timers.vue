<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
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

import { useUserAccessFlagChecker } from '@/api/index.js';
import { EditableTimer, useTimersApi, type TimerResponse } from '@/api/timers.js';
import Modal from '@/components/timers/modal.vue';
import { renderIcon } from '@/helpers/index.js';

const { t } = useI18n();
const userCanManageTimers = useUserAccessFlagChecker('MANAGE_TIMERS');

const timersApi = useTimersApi();
const timersRemove = timersApi.useMutationRemoveTimer();
const timersUpdate = timersApi.useMutationUpdateTimer();

const { data, fetching } = timersApi.useQueryTimers();
const timers = computed(() => {
	return data.value?.timers ?? [];
});

const columns = computed<DataTableColumns<TimerResponse>>(() => [
	{
		title: t('sharedTexts.name'),
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
	{
		title: t('timers.table.columns.intervalInMessages'),
		key: 'messageInterval',
		render(row) {
			return h(NTag, { type: 'info' }, { default: () => `${row.messageInterval}` });
		},
	},
	{
		title: t('sharedTexts.status'),
		key: 'enabled',
		render(row) {
			return h(NSwitch, {
				value: row.enabled,
				disabled: !userCanManageTimers.value,
				onUpdateValue: (enabled) => {
					timersUpdate.executeMutation({ id: row.id, opts: { enabled } });
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
					onPositiveClick: () => timersRemove.executeMutation({ id: row.id }),
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

function openModal(t: EditableTimer | null) {
	editableTimer.value = t;
	showModal.value = true;
}

function closeModal() {
	showModal.value = false;
}

const timersLength = computed(() => timers.value.length);
</script>

<template>
	<n-space justify="space-between" align="center" class="mb-2">
		<h1 class="text-2xl">
			{{ t('sidebar.timers') }}
		</h1>
		<div>
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
		:isLoading="fetching"
		:columns="columns"
		:data="timers"
	/>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableTimer?.name ?? t('timers.newTimer')"
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
