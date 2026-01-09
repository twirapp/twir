<script setup lang='ts'>
import { PlusIcon } from 'lucide-vue-next'
import { ref } from 'vue'


import RoleModal from './ui/modal.vue'

import type { ChannelRolesQuery } from '~/gql/graphql'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useRoles } from '#layers/dashboard/api/roles'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'




import { ChannelRolePermissionEnum, RoleTypeEnum } from '~/gql/graphql'

const rolesManager = useRoles()
const { data: roles } = rolesManager.useRolesQuery()
const rolesDeleter = rolesManager.useRolesDeleteMutation()

const editableRole = ref<ChannelRolesQuery['roles'][number] | null>(null)
const showModal = ref(false)
const showDelete = ref(false)
const roleToDelete = ref<string | null>(null)

function openModal(role: ChannelRolesQuery['roles'][number] | null) {
	editableRole.value = role
	showModal.value = true
}

function openDelete(roleId: string) {
	roleToDelete.value = roleId
	showDelete.value = true
}

async function deleteRole() {
	if (roleToDelete.value) {
		await rolesDeleter.executeMutation({ id: roleToDelete.value })
	}
}

const userCanManageRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageRoles)

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-2">
		<UiCard
			class="min-w-[400px]"
			:class="{ 'cursor-pointer': userCanManageRoles, 'cursor-not-allowed': !userCanManageRoles }"
			@click="() => {
				if (userCanManageRoles) {
					openModal(null)
				}
			}"
		>
			<UiCardContent class="flex items-center justify-center p-6">
				<PlusIcon class="h-8 w-8" />
			</UiCardContent>
		</UiCard>

		<UiCard
			v-for="role in roles?.roles"
			:key="role.id"
			class="min-w-[400px]"
		>
			<UiCardContent class="flex items-center justify-between p-6">
				<span class="text-2xl">
					{{ role.name }}
				</span>
				<div class="flex gap-2">
					<UiButton
						:disabled="!userCanManageRoles"
						variant="outline"
						@click="openModal(role)"
					>
						{{ t('sharedButtons.edit') }}
					</UiButton>

					<UiButton
						v-if="role.type === RoleTypeEnum.Custom"
						:disabled="role.type !== 'CUSTOM' || !userCanManageRoles"
						variant="destructive"
						@click="openDelete(role.id)"
					>
						{{ t('sharedButtons.delete') }}
					</UiButton>
				</div>
			</UiCardContent>
		</UiCard>
	</div>

	<UiDialog v-model:open="showModal">
		<DialogOrSheet>
			<UiDialogHeader>
				<UiDialogTitle>
					{{ editableRole?.name || 'Create role' }}
				</UiDialogTitle>
			</UiDialogHeader>

			<RoleModal :role="editableRole" />
		</DialogOrSheet>
	</UiDialog>

	<UiActionConfirm
		v-model:open="showDelete"
		@confirm="deleteRole"
	/>
</template>
