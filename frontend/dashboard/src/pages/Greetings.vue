<script setup lang='ts'>
import { IconPencil, IconTrash } from '@tabler/icons-vue';
import type { Greeting } from '@twir/grpc/generated/api/api/greetings';
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
import { h, ref, watch, onMounted, computed } from 'vue';

import { useGreetingsManager, useTwitchGetUsers } from '@/api/index.js';
import Modal from '@/components/greetings/modal.vue';
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

const columns: DataTableColumns<Greeting> = [
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
		title: 'User name',
		key: 'userName',
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
		title: 'Text',
		key: 'text',
		render(row) {
			return h(NTag, { type: 'info', bordered: true }, { default: () => row.text });
		},
	},
	{
		title: 'Status',
		key: 'enabled',
		width: 100,
		render(row) {
			return h(
				NSwitch,
				{
					value: row.enabled,
					onUpdateValue: (value: boolean) => {
						greetingsPatcher.mutateAsync({ id: row.id, enabled: value }).then(() => {
							row.enabled = value;
						});
					},
				},
				{ default: () => row.enabled },
			);
		},
	},
	{
		title: 'Actions',
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
				}, {
					icon: renderIcon(IconPencil),
				});

			const deleteButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => greetingsDeleter.mutate({ id: row.id! }),
				},
				{
					trigger: () => h(NButton, {
						type: 'error',
						size: 'small',
						quaternary: true,
					}, {
						default: renderIcon(IconTrash),
					}),
					default: () => 'Are you sure you want to delete this variable?',
				},
			);

			return h(NSpace, { }, { default: () => [editButton, deleteButton] });
		},
	},
];

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
    <h2>Greetings</h2>
    <n-button secondary type="success" @click="openModal(null)">
      Create
    </n-button>
  </n-space>
  <n-alert>
    <p>Greeting system used for welcoming new users typed their first message on stream.</p>
    <p>
      If you wanna greet every user in chat, not only specified - then you can use events system.
    </p>
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
    :title="editableGreeting?.userName || 'Create greeting'"
    class="modal"
    :style="{
      width: '400px',
      top: '50px',
    }"
    :on-close="closeModal"
  >
    <modal :greeting="editableGreeting" @close="closeModal" />
  </n-modal>
</template>
