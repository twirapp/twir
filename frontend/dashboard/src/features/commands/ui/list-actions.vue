<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

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
			enabled: newValue,
		},
	})

	toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
		duration: 1500,
	})
}

async function deleteCommand() {
	await deleter.executeMutation({ id: props.row.id })

	toast({
		title: t('sharedTexts.deleted'),
		variant: 'success',
		duration: 1500,
	})
}
</script>

<template>
	<div class="flex items-center gap-4">
		<Switch
			:disabled="!userCanManageCommands"
			:checked="row.enabled"
			@update:checked="switchEnabled"
		/>
		<div class="flex gap-2">
			<RouterLink v-slot="{ href, navigate }" custom :to="`/dashboard/commands/custom/${row.id}`">
				<Button
					as="a"
					:href="href"
					:disabled="!userCanManageCommands"
					variant="secondary"
					size="icon"
					@click="navigate"
				>
					<PencilIcon class="h-4 w-4" />
				</Button>
			</RouterLink>
			<Button
				v-if="row.module === 'CUSTOM'"
				:disabled="!userCanManageCommands"
				variant="destructive"
				size="icon"
				@click="showDelete = true"
			>
				<TrashIcon class="h-4 w-4" />
			</Button>
		</div>
	</div>

	<ActionConfirmation
		v-model:open="showDelete"
		@confirm="deleteCommand"
	/>
</template>
