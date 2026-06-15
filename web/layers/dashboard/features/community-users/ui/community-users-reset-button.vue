<script setup lang="ts">
import { ref } from 'vue'
import { useCommunityUsersApi } from '~~/layers/dashboard/api/community-users.js'

import ActionConfirm from '@/components/ui/action-confirm'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { CommunityUsersResetType } from '~/gql/graphql.js'

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

const selectOptions = Object.values(CommunityUsersResetType).map((v) => {
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
			<Button
				variant="destructive"
				size="sm"
			>
				<Icon
					name="lucide:trash"
					class="mr-2 h-3.5 w-3.5 text-white"
				/>
				{{ t('community.users.reset.label') }}
				<Icon
					name="lucide:chevron-down"
					class="ml-2 h-3.5 w-3.5 text-white"
				/>
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
				@click="
					() => {
						selectedType = option.value
						showConfirm = true
					}
				"
			>
				{{ option.label }}
			</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>

	<ActionConfirm
		v-model:open="showConfirm"
		:confirm-text="
			t('community.users.reset.resetQuestion', {
				title: `${selectedType.toLowerCase().split('_').join(' ')}`,
			})
		"
		@confirm="resetColumn"
	/>
</template>
