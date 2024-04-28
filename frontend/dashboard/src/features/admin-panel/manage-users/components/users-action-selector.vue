<script setup lang="ts">
import { BanIcon, LogOutIcon, SwordIcon, WrenchIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useUsers } from '../composables/use-users.js'

import { useAdminActions } from '@/api/admin/actions.js'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'

defineProps<{
	userId: string
	isBanned: boolean
	isBotAdmin: boolean
}>()

const { t } = useI18n()
const { switchBan, switchAdmin } = useUsers()

const adminActions = useAdminActions()
const dropAuthSession = adminActions.useMutationDropUserAuthSession()
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger as-child>
			<Button variant="secondary" size="icon">
				<WrenchIcon class="size-4" />
			</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent align="end">
			<DropdownMenuItem @click="dropAuthSession.executeMutation({ userId })">
				<LogOutIcon class="mr-2 h-4 w-4" />
				<span>{{ t('adminPanel.manageUsers.dropSession') }}</span>
			</DropdownMenuItem>

			<DropdownMenuItem @click="switchAdmin.executeMutation({ userId })">
				<SwordIcon class="mr-2 h-4 w-4" />
				<span>{{ isBotAdmin ? t('adminPanel.manageUsers.unMod') : t('adminPanel.manageUsers.giveMod') }}</span>
			</DropdownMenuItem>

			<DropdownMenuItem @click="switchBan.executeMutation({ userId })">
				<BanIcon class="mr-2 h-4 w-4" />
				<span>{{ isBanned ? t('adminPanel.manageUsers.unBan') : t('adminPanel.manageUsers.giveBan') }}</span>
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
