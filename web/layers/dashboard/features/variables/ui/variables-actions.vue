<script lang="ts" setup>
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { type CustomVariable, useVariablesApi } from '#layers/dashboard/api/variables'


import { toast } from 'vue-sonner'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const props = defineProps<{ row: CustomVariable }>()

const manager = useVariablesApi()
const deleter = manager.useMutationRemoveVariable()

const userCanManageVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageVariables)
const { t } = useI18n()

const showDelete = ref(false)

async function deleteVariable() {
	await deleter.executeMutation({ id: props.row.id })

	toast.success(t('sharedTexts.deleted'), {
		duration: 1500,
	})
}
</script>

<template>
	<div class="flex gap-2 items-center justify-end">
		<RouterLink v-slot="{ href, navigate }" custom :to="`/dashboard/variables/${row.id}`">
			<UiButton
				as="a"
				:href="href"
				:disabled="!userCanManageVariables"
				variant="secondary"
				size="icon"
				@click="navigate"
			>
				<PencilIcon class="h-4 w-4" />
			</UiButton>
		</RouterLink>
		<UiButton
			:disabled="!userCanManageVariables"
			variant="destructive"
			size="icon"
			@click="showDelete = true"
		>
			<TrashIcon class="h-4 w-4" />
		</UiButton>
	</div>

	<UiActionConfirm v-model:open="showDelete" @confirm="deleteVariable" />
</template>
