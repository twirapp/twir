<script setup lang="ts">
import { ChevronDownIcon, Trash } from 'lucide-vue-next'
import { ref } from 'vue'


import { useCommunityUsersApi } from '#layers/dashboard/api/community-users.ts'



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

const selectOptions = Object.values(CommunityUsersResetType).map(v => {
	const label = v.toLowerCase().split('_').join(' ')

	return {
		label: label.charAt(0).toUpperCase() + label.slice(1),
		value: v,
	}
})
</script>

<template>
	<UiDropdownMenu>
		<UiDropdownMenuTrigger as-child>
			<UiButton variant="destructive" size="sm">
				<Trash class="mr-2 h-3.5 w-3.5 text-white" />
				{{ t('community.users.reset.label') }}
				<ChevronDownIcon class="ml-2 h-3.5 w-3.5 text-white" />
			</UiButton>
		</UiDropdownMenuTrigger>
		<UiDropdownMenuContent>
			<UiDropdownMenuLabel>
				{{ t('community.users.reset.label') }}
			</UiDropdownMenuLabel>
			<UiDropdownMenuSeparator />
			<UiDropdownMenuItem
				v-for="option of selectOptions"
				:key="option.value"
				@click="() => {
					selectedType = option.value
					showConfirm = true
				}"
			>
				{{ option.label }}
			</UiDropdownMenuItem>
		</UiDropdownMenuContent>
	</UiDropdownMenu>

	<UiActionConfirm
		v-model:open="showConfirm"
		:confirm-text="t('community.users.reset.resetQuestion', { title: `${selectedType.toLowerCase().split('_').join(' ')}` })"
		@confirm="resetColumn"
	/>
</template>
