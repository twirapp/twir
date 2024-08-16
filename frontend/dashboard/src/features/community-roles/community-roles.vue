<script setup lang='ts'>
import { IconPlus } from '@tabler/icons-vue'
import {
	NButton,
	NCard,
	NModal,
	NPopconfirm,
	NSpace,
	NText,
} from 'naive-ui'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import RoleModal from './ui/modal.vue'

import type { ChannelRolesQuery } from '@/gql/graphql'

import { useUserAccessFlagChecker } from '@/api/index.js'
import { useRoles } from '@/api/roles'
import { ChannelRolePermissionEnum, RoleTypeEnum } from '@/gql/graphql'

const rolesManager = useRoles()
const { data: roles } = rolesManager.useRolesQuery()
const rolesDeleter = rolesManager.useRolesDeleteMutation()

const editableRole = ref<ChannelRolesQuery['roles'][number] | null>(null)
const showModal = ref(false)
function openModal(role: ChannelRolesQuery['roles'][number] | null) {
	editableRole.value = role
	showModal.value = true
}
const closeModal = () => showModal.value = false

const userCanManageRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageRoles)

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-2">
		<NCard
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
			<NSpace align="center" justify="center" vertical>
				<NText class="text-[30px]">
					<IconPlus />
				</NText>
			</NSpace>
		</NCard>
		<NCard
			v-for="role in roles?.roles"
			:key="role.id"
			size="small"
			class="min-w-[400px]"
			hoverable
		>
			<NSpace justify="space-between" align="center">
				<NText class="text-[30px]">
					{{ role.name }}
				</NText>
				<NSpace>
					<NButton :disabled="!userCanManageRoles" secondary type="success" @click="openModal(role)">
						{{ t('sharedButtons.edit') }}
					</NButton>
					<NPopconfirm
						v-if="role.type === RoleTypeEnum.Custom"
						:positive-text="t('deleteConfirmation.confirm')"
						:negative-text="t('deleteConfirmation.cancel')"
						@positive-click="() => rolesDeleter.executeMutation({ id: role.id })"
					>
						<template #trigger>
							<NButton :disabled="role.type !== 'CUSTOM' || !userCanManageRoles" secondary type="error">
								{{ t('sharedButtons.delete') }}
							</NButton>
						</template>
						{{ t('deleteConfirmation.text') }}
					</NPopconfirm>
				</NSpace>
			</NSpace>
		</NCard>

		<NModal
			v-model:show="showModal"
			:mask-closable="false"
			:segmented="true"
			preset="card"
			:title="editableRole?.name || 'Create role'"
			:style="{ width: '600px', top: '50px' }"
			:on-close="closeModal"
		>
			<RoleModal :role="editableRole" @close="closeModal" />
		</NModal>
	</div>
</template>
