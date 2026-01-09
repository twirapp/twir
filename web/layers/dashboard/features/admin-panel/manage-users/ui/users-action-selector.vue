<script setup lang="ts">
import { BanIcon, SwordIcon, WrenchIcon } from 'lucide-vue-next'


import { useUsers } from '../composables/use-users.js'

// ...existing code...

defineProps<{
	userId: string
	isBanned: boolean
	isBotAdmin: boolean
}>()

const { t } = useI18n()
const { switchBan, switchAdmin } = useUsers()
const { user: profile } = storeToRefs(useDashboardAuth())
</script>

<template>
	<UiDropdownMenu v-if="userId !== profile?.id">
		<UiDropdownMenuTrigger as-child>
			<UiButton variant="secondary" size="icon">
				<WrenchIcon class="size-4" />
			</UiButton>
		</UiDropdownMenuTrigger>
		<UiDropdownMenuContent align="end">
			<UiDropdownMenuItem @click="switchAdmin.executeMutation({ userId })">
				<SwordIcon class="mr-2 h-4 w-4" />
				<span>{{ isBotAdmin ? t('adminPanel.manageUsers.unMod') : t('adminPanel.manageUsers.giveMod') }}</span>
			</UiDropdownMenuItem>

			<UiDropdownMenuItem @click="switchBan.executeMutation({ userId })">
				<BanIcon class="mr-2 h-4 w-4" />
				<span>{{ isBanned ? t('adminPanel.manageUsers.unBan') : t('adminPanel.manageUsers.giveBan') }}</span>
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>
</template>
