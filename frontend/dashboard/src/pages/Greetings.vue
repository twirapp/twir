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

import { useGreetingsManager, useTwitchGetUsers, useUserAccessFlagChecker } from '@/api/index.js';
import GreetingsModal from '@/components/greetings/modal.vue';
import { EditableGreeting } from '@/components/greetings/types.js';
import { renderIcon } from '@/helpers/index.js';

const greetingsManager = useGreetingsManager();
const greetings = greetingsManager.getAll({});
const greetingsDeleter = greetingsManager.deleteOne;
const greetingsPatcher = greetingsManager.patch!;

const showModal = ref(false);

const twitchUsersIds = computed(() => {
	return greetings.data.value?.greetings.map((g) => g.userId) ?? [];
});
const twitchUsers = useTwitchGetUsers({
	ids: twitchUsersIds,
});

const userCanManageGreetings = useUserAccessFlagChecker('MANAGE_GREETINGS');

const { t } = useI18n();

const columns = computed<DataTableColumns<EditableGreeting>>(() => [
	{
		title: '',
		key: 'avatar',
		width: 50,
		render(row) {
			const user = twitchUsers.data.value?.users.find((u) => u.id === row.userId);
			if (!user) return;

			return h(NAvatar, { size: 'medium', src: user.profileImageUrl, round: true });
		},
	},
	{
		title: t('sharedTexts.userName'),
		key: 'userName',
		maxWidth: 150,
		render(row) {
			return h(NTag, { type: 'info', bordered: false }, {
				default: () => {
					const user = twitchUsers.data.value?.users.find((u) => u.id === row.userId);
					return user ? user.displayName : 'Unknown';
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
					onUpdateValue: (value: boolean) => {
						greetingsPatcher.mutate({ id: row.id!, enabled: value });
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
					onPositiveClick: () => greetingsDeleter.mutate({ id: row.id! }),
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

const editableGreeting = ref<EditableGreeting | null>(null);

function openModal(t: EditableGreeting | null) {
	const twitchUser = twitchUsers.data.value?.users.find((u) => u.id === t?.userId);
	editableGreeting.value = t ? {
		...t,
		userName: twitchUser?.login || 'Unknown user',
	} : null;

	showModal.value = true;
}

function closeModal() {
	showModal.value = false;
}
</script>

<template>
	<n-space justify="space-between" align="center">
		<h2>{{ t('greetings.title') }}</h2>
		<n-button :disabled="!userCanManageGreetings" secondary type="success" @click="openModal(null)">
			{{ t('sharedButtons.create') }}
		</n-button>
	</n-space>
	<n-alert type="info" :title="t('greetings.info.title')">
		{{ t('greetings.info.text') }}
	</n-alert>
	<n-data-table
		:isLoading="greetings.isLoading.value || twitchUsers.isLoading.value"
		:columns="columns"
		:data="greetings.data.value?.greetings ?? []"
		style="margin-top: 20px;"
	/>

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableGreeting?.userName || 'Create'"
		class="modal"
		:style="{
			width: '400px',
			top: '50px',
		}"
		:on-close="closeModal"
	>
		<greetings-modal :greeting="editableGreeting" @close="closeModal" />
	</n-modal>
</template>
