<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useUsers } from '../composables/use-users.js'

import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

defineProps<{
	userId: string
	isBanned: boolean
	isBotAdmin: boolean
}>()

const { t } = useI18n()
const { switchBan, switchAdmin } = useUsers()
const { data: profile } = useProfile()
</script>

<template>
	<DropdownMenu v-if="userId !== profile?.id">
		<DropdownMenuTrigger as-child>
			<Button variant="secondary" size="icon">
				<Icon name="lucide:wrench" class="size-4" />
			</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent align="end">
			<DropdownMenuItem @click="switchAdmin.executeMutation({ userId })">
				<Icon name="lucide:sword" class="mr-2 h-4 w-4" />
				<span>{{ isBotAdmin ? t('adminPanel.manageUsers.unMod') : t('adminPanel.manageUsers.giveMod') }}</span>
			</DropdownMenuItem>

			<DropdownMenuItem @click="switchBan.executeMutation({ userId })">
				<Icon name="lucide:ban" class="mr-2 h-4 w-4" />
				<span>{{ isBanned ? t('adminPanel.manageUsers.unBan') : t('adminPanel.manageUsers.giveBan') }}</span>
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</template>
