<script setup lang='ts'>
import { IconPlus } from '@tabler/icons-vue';
import {
  NCard,
  NSpace,
  NText,
  NModal,
  NButton,
  NPopconfirm,
	NResult,
} from 'naive-ui';
import { ref } from 'vue';

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

const userCanViewRoles = useUserAccessFlagChecker('VIEW_ROLES');
const userCanManageRoles = useUserAccessFlagChecker('MANAGE_ROLES');
</script>

<template>
	<n-result v-if="!userCanViewRoles" status="403" title="You haven't acces to view roles" />
	<div v-else>
		<n-space align="center" justify="center" vertical>
			<n-card
				class="card"
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
					<n-text class="text">
						<IconPlus />
					</n-text>
				</n-space>
			</n-card>
			<n-card
				v-for="role in roles?.roles"
				:key="role.id"
				size="small"
				class="card"
				hoverable
			>
				<n-space justify="space-between" align="center">
					<n-text class="text">
						{{ role.name }}
					</n-text>
					<n-space>
						<n-button :disabled="!userCanManageRoles" secondary type="success" @click="openModal(role)">
							Edit
						</n-button>
						<n-popconfirm @positive-click="() => rolesDeleter.mutateAsync({ id: role.id })">
							<template #trigger>
								<n-button :disabled="role.type !== 'CUSTOM' || !userCanManageRoles" secondary type="error">
									Remove
								</n-button>
							</template>
							Are you sure?
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
	</div>
</template>

<style scoped>
.card {
	min-width: 400px;
}

.card .text {
	font-size: 30px;
}
</style>
