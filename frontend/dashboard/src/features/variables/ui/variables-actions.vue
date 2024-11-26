<script lang="ts" setup>
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api'
import { type CustomVariable, useVariablesApi } from '@/api/variables'
import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{ row: CustomVariable }>()

const manager = useVariablesApi()
const deleter = manager.useMutationRemoveVariable()

const userCanManageVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageVariables)
const { t } = useI18n()
const { toast } = useToast()

const showDelete = ref(false)

async function deleteVariable() {
	await deleter.executeMutation({ id: props.row.id })

	toast({
		title: t('sharedTexts.deleted'),
		variant: 'success',
		duration: 1500,
	})
}
</script>

<template>
	<div class="flex gap-2 items-center justify-end">
		<RouterLink v-slot="{ href, navigate }" custom :to="`/dashboard/variables/${row.id}`">
			<Button
				as="a"
				:href="href"
				:disabled="!userCanManageVariables"
				variant="secondary"
				size="icon"
				@click="navigate"
			>
				<PencilIcon class="h-4 w-4" />
			</Button>
		</RouterLink>
		<Button
			:disabled="!userCanManageVariables"
			variant="destructive"
			size="icon"
			@click="showDelete = true"
		>
			<TrashIcon class="h-4 w-4" />
		</Button>
	</div>

	<ActionConfirm v-model:open="showDelete" @confirm="deleteVariable" />
</template>
