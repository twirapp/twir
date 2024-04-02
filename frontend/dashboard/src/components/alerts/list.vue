<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import { Alert } from '@twir/api/messages/alerts/alerts';
import {
	DataTableColumns,
	NButton,
	NDataTable,
	NModal,
	NPopconfirm,
	NSpace,
	NTag,
	NImage,
} from 'naive-ui';
import { computed, h, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	useAlertsManager,
	useCommandsManager,
	useTwitchRewards,
	useUserAccessFlagChecker,
} from '@/api';
import AlertModal from '@/components/alerts/modal.vue';
import { type EditableAlert } from '@/components/alerts/types.js';
import { renderIcon } from '@/helpers';

const props = withDefaults(defineProps<{
	withSelect: boolean
}>(), {
	withSelect: false,
});

const emits = defineEmits<{
	select: [id: string]
	delete: [id: string]
}>();

const manager = useAlertsManager();
const deleter = manager.deleteOne;
const { data, isLoading } = manager.getAll({});

const { t } = useI18n();

const userCanManageAlerts = useUserAccessFlagChecker('MANAGE_ALERTS');
const { data: rewards } = useTwitchRewards();
const { data: commands } = useCommandsManager().getAll({});

const columns = computed<DataTableColumns<Alert>>(() => [
	{
		title: t('alerts.name'),
		key: 'text',
		render(row) {
			return h(
				NTag,
				{ type: 'info', bordered: false },
				{
					default: () => row.name,
				},
			);
		},
	},
	{
		title: 'Rewards',
		key: 'rewardId',
		render(row) {
			const selectedRewards = rewards?.value?.rewards.filter(r => row.rewardIds.includes(r.id));
			if (!selectedRewards?.length) {
				return '';
			}

			const mappedRewards = selectedRewards.map(r => h(NSpace, {
				type: 'success',
				bordered: false,
				style: `background-color: ${r.backgroundColor}`,
				class: 'rounded-lg p-1',
				align: 'center',
			}, {
				default: () => [
					h(NImage, {
						src: r.image?.url1X || r.defaultImage?.url4X,
						class: 'flex items-center h-5 w-5',
					}),
					r.title,
				],
			}));

			return h(
				NSpace,
				{ vertical: true },
				{
					default: () => mappedRewards,
				},
			);
		},
	},
	{
		title: 'Commands',
		key: 'commands',
		render(row) {
			const selectedCommands = commands?.value?.commands.filter(r => row.commandIds.includes(r.id));
			if (!selectedCommands?.length) {
				return '';
			}

			return h(
				NSpace,
				{ vertical: true },
				{
					default: () => selectedCommands.map(c => `!${c.name}`),
				},
			);
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
					disabled: !userCanManageAlerts.value,
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => deleter.mutate({ id: row.id }),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
						disabled: !userCanManageAlerts.value,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => t('deleteConfirmation.text'),
				},
			);

			const selectButton = h(
				NButton,
				{
					type: 'success',
					size: 'small',
					block: true,
					onClick: () => emits('select', row.id),
					secondary: true,
					disabled: !userCanManageAlerts.value,
				}, {
					default: () => t('sharedButtons.select'),
				});

			const buttons = [editButton, deleteButton];

			if (props.withSelect) {
				buttons.unshift(selectButton);
			}

			return h(NSpace, {}, { default: () => buttons });
		},
	},
]);


const showModal = ref(false);
const editableAlert = ref<EditableAlert | null>(null);

function openModal(t: EditableAlert | null) {
	editableAlert.value = t;
	showModal.value = true;
}
</script>

<template>
	<n-space justify="space-between" align="center">
		<h2>{{ t('alerts.title') }}</h2>
		<n-button :disabled="!userCanManageAlerts" secondary type="success" @click="openModal(null)">
			{{ t('sharedButtons.create') }}
		</n-button>
	</n-space>

	<n-data-table
		:isLoading="isLoading"
		:columns="columns"
		:data="data?.alerts ?? []"
	/>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableAlert?.name"
		class="modal"
		:style="{
			width: '400px',
			top: '50px',
		}"
		:on-close="() => showModal = false"
	>
		<alert-modal :alert="editableAlert" @close="() => showModal = false" />
	</n-modal>
</template>
