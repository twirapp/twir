<script setup lang='ts'>
import { IconPlus } from '@tabler/icons-vue';
import {
  NCard,
  NSpace,
  NText,
  NModal,
  NButton,
  NPopconfirm,
} from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useRolesManager, useUserAccessFlagChecker } from '@/api/index.js';
import RoleModal from '@/components/roles/modal.vue';
import type { EditableRole } from '@/components/roles/types.js';

const rolesManager = useRolesManager();
const { data: roles } = rolesManager.getAll({});
const rolesDeleter = rolesManager.deleteOne;

const editableRole = ref<EditableRole | null>(null);
const showModal = ref(false);
function openModal(role: EditableRole | null) {
	editableRole.value = role;
	showModal.value = true;
}
const closeModal = () => showModal.value = false;

const userCanManageRoles = useUserAccessFlagChecker('MANAGE_ROLES');

const { t } = useI18n();
</script>

<template>
	<n-space align="center" justify="center" vertical>
		<n-card
			class="min-w-[400px]"
			:style="{ cursor: userCanManageRoles ? 'pointer' : 'not-allowed' }"
			size="small"
			bordered
			hoverable
			@click="() => {
				if (userCanManageRoles) {
					openModal(null)
				}
			}"
		>
			<n-space align="center" justify="center" vertical>
				<n-text class="text-[30px]">
					<IconPlus />
				</n-text>
			</n-space>
		</n-card>
		<n-card
			v-for="role in roles?.roles"
			:key="role.id"
			size="small"
			class="min-w-[400px]"
			hoverable
		>
			<n-space justify="space-between" align="center">
				<n-text class="text-[30px]">
					{{ role.name }}
				</n-text>
				<n-space>
					<n-button :disabled="!userCanManageRoles" secondary type="success" @click="openModal(role)">
						{{ t('sharedButtons.edit') }}
					</n-button>
					<n-popconfirm
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => rolesDeleter.mutateAsync({ id: role.id })"
					>
						<template #trigger>
							<n-button :disabled="role.type !== 'CUSTOM' || !userCanManageRoles" secondary type="error">
								{{ t('sharedButtons.delete') }}
							</n-button>
						</template>
						{{ t('deleteConfirmation.text') }}
					</n-popconfirm>
				</n-space>
			</n-space>
		</n-card>

		<n-modal
			v-model:show="showModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="editableRole?.name || 'Create role'"
			:style="{ width: '600px',top: '50px' }"
			:on-close="closeModal"
		>
			<role-modal :role="editableRole" @close="closeModal" />
		</n-modal>
	</n-space>
</template>
