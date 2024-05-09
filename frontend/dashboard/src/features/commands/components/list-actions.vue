<script setup lang="ts">
import { IconPencil, IconTrash } from '@tabler/icons-vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommandEdit } from '../composables/use-command-edit'

import type { Command } from '@/gql/graphql'

import { useUserAccessFlagChecker } from '@/api'
import { useCommandsApi } from '@/api/commands/commands.js'
import ActionConfirmation from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/toast/use-toast'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{ row: Command }>()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)

const manager = useCommandsApi()
const deleter = manager.useMutationDeleteCommand()
const patcher = manager.useMutationUpdateCommand()!

const { t } = useI18n()
const { toast } = useToast()

const showDelete = ref(false)

async function switchEnabled(newValue: boolean) {
	await patcher?.executeMutation({
		id: props.row.id,
		opts: {
			enabled: newValue
		}
	})

	toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
		duration: 1500
	})
}

async function deleteCommand() {
	await deleter.executeMutation({ id: props.row.id })

	toast({
		title: t('sharedTexts.deleted'),
		variant: 'success',
		duration: 1500
	})
}

const commandEdit = useCommandEdit()
</script>

<template>
	<div class="flex items-center gap-4">
		<Switch
			:disabled="!userCanManageCommands"
			:checked="row.enabled"
			@update:checked="switchEnabled"
		/>
		<div class="flex gap-0.5">
			<Button
				v-if="row.module === 'CUSTOM'"
				:disabled="!userCanManageCommands"
				variant="ghost"
				size="sm"
				@click="showDelete = true"
			>
				<IconTrash class="h-5 w-5" />
			</Button>
			<Button
				:disabled="!userCanManageCommands"
				size="sm"
				variant="ghost"
				@click="commandEdit.editCommand(row.id)"
			>
				<IconPencil class="h-5 w-5" />
			</Button>
		</div>
	</div>

	<ActionConfirmation
		v-model:open="showDelete"
		@confirm="deleteCommand"
	/>
</template>
