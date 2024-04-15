<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import {
	type DataTableColumns,
	NDataTable,
	NSpace,
	NTag,
	NAlert,
	NButton,
	NPopconfirm,
	NModal,
	NSwitch,
	NAvatar,
} from 'naive-ui';
import { h, ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api/index.js';
import GreetingsModal from '@/components/greetings/modal.vue';
import { renderIcon } from '@/helpers/index.js';
import { useGreetingsApi, type Greetings } from '@/api/greetings.js';

const { t } = useI18n();
const userCanManageGreetings = useUserAccessFlagChecker('MANAGE_GREETINGS');
const showModal = ref(false);

const greetingsApi = useGreetingsApi();
const greetingsUpdate = greetingsApi.useMutationUpdateGreetings();
const greetingsRemove = greetingsApi.useMutationRemoveGreetings();
const { data: greetingsData, fetching } = greetingsApi.useQueryGreetings();

const greetings = computed(() => {
	return greetingsData.value?.greetings ?? [];
})

const columns = computed<DataTableColumns<Greetings>>(() => [
	{
		title: '',
		key: 'avatar',
		width: 50,
		render(row) {
			if (!row.twitchProfile) return;

			return h(NAvatar, { size: 'medium', src: row.twitchProfile.profileImageUrl, round: true });
		},
	},
	{
		title: t('sharedTexts.userName'),
		key: 'userName',
		maxWidth: 150,
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, {
				default: () => {
					return row.twitchProfile ? row.twitchProfile.displayName : 'Unknown';
				},
			});
		},
	},
	{
		title: t('sharedTexts.response'),
		key: 'text',
		maxWidth: 600,
		render(row) {
			return row.text;
		},
	},
	{
		title: t('sharedTexts.status'),
		key: 'enabled',
		width: 100,
		render(row) {
			return h(
				NSwitch,
				{
					value: row.enabled,
					onUpdateValue: (enabled) => {
						greetingsUpdate.executeMutation({ id: row.id, opts: { enabled } });
					},
					disabled: !userCanManageGreetings.value,
				},
				{ default: () => row.enabled },
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
					disabled: !userCanManageGreetings.value,
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => greetingsRemove.executeMutation({ id: row.id }),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
						disabled: !userCanManageGreetings.value,
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

const editableGreeting = ref<Greetings | null>(null);

function openModal(greetings: Greetings | null) {
	editableGreeting.value = greetings
	showModal.value = true;
}

function closeModal() {
	showModal.value = false;
}
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('greetings.title') }}
			</h4>
			<n-button :disabled="!userCanManageGreetings" secondary type="success" @click="openModal(null)">
				{{ t('sharedButtons.create') }}
			</n-button>
		</div>

		<n-alert type="info" :title="t('greetings.info.title')">
			{{ t('greetings.info.text') }}
		</n-alert>

		<n-data-table
			:isLoading="fetching"
			:columns="columns"
			:data="greetings"
		/>
	</div>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableGreeting ? t('greetings.edit') : t('greetings.create')"
		class="modal"
		:style="{ width: '400px' }"
		:on-close="closeModal"
	>
		<greetings-modal :greeting="editableGreeting" @close="closeModal" />
	</n-modal>
</template>
