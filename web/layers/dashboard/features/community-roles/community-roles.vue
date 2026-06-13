<script setup lang='ts'>
import { PlusIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import RoleModal from './ui/modal.vue'

import type { ChannelRolesQuery } from '@/gql/graphql'

import { useUserAccessFlagChecker } from '@/api/auth'
import { useRoles } from '@/api/roles'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Dialog, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { ChannelRolePermissionEnum, RoleTypeEnum } from '@/gql/graphql'

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
		<Card
			class="min-w-[400px]"
			:class="{ 'cursor-pointer': userCanManageRoles, 'cursor-not-allowed': !userCanManageRoles }"
			@click="() => {
				if (userCanManageRoles) {
					openModal(null)
				}
			}"
		>
			<CardContent class="flex items-center justify-center p-6">
				<PlusIcon class="h-8 w-8" />
			</CardContent>
		</Card>

		<Card
			v-for="role in roles?.roles"
			:key="role.id"
			class="min-w-[400px]"
		>
			<CardContent class="flex items-center justify-between p-6">
				<span class="text-2xl">
					{{ role.name }}
				</span>
				<div class="flex gap-2">
					<Button
						:disabled="!userCanManageRoles"
						variant="outline"
						@click="openModal(role)"
					>
						{{ t('sharedButtons.edit') }}
					</Button>

					<Button
						v-if="role.type === RoleTypeEnum.Custom"
						:disabled="role.type !== 'CUSTOM' || !userCanManageRoles"
						variant="destructive"
						@click="openDelete(role.id)"
					>
						{{ t('sharedButtons.delete') }}
					</Button>
				</div>
			</CardContent>
		</Card>
	</div>

	<Dialog v-model:open="showModal">
		<DialogOrSheet>
			<DialogHeader>
				<DialogTitle>
					{{ editableRole?.name || 'Create role' }}
				</DialogTitle>
			</DialogHeader>

			<RoleModal :role="editableRole" />
		</DialogOrSheet>
	</Dialog>

	<ActionConfirm
		v-model:open="showDelete"
		@confirm="deleteRole"
	/>
</template>
