<script setup lang="ts">
import { Trash } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommunityUsersApi } from '@/api/community-users.ts'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { CommunityUsersResetType } from '@/gql/graphql.js'

const communityUsersApi = useCommunityUsersApi()
const communityResetMutation = communityUsersApi.useMutationCommunityReset()

const { t } = useI18n()
const selectedType = ref(CommunityUsersResetType.Watched)

const showConfirm = ref(false)
async function resetColumn() {
	await communityResetMutation.executeMutation({
		type: selectedType.value,
	})
}

const selectOptions = Object.values(CommunityUsersResetType).map(v => {
	const label = v.toLowerCase().split('_').join(' ')

	return {
		label: label.charAt(0).toUpperCase() + label.slice(1),
		value: v,
	}
})
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger as-child>
			<Button variant="destructive" size="sm">
				<Trash class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
				{{ t('community.users.reset.label') }}
			</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent>
			<DropdownMenuLabel>
				{{ t('community.users.reset.label') }}
			</DropdownMenuLabel>
			<DropdownMenuSeparator />
			<DropdownMenuItem
				v-for="option of selectOptions"
				:key="option.value"
				@click="() => {
					selectedType = option.value
					showConfirm = true
				}"
			>
				{{ option.label }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>

	<ActionConfirm
		v-model:open="showConfirm"
		:confirm-text="t('community.users.reset.resetQuestion', { title: `${selectedType.toLowerCase().split('_').join(' ')}` })"
		@confirm="resetColumn"
	/>
</template>
