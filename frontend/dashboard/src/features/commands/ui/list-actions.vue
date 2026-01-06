<script setup lang="ts">
import { CopyIcon, PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import type { Command } from '@/gql/graphql'

import { useUserAccessFlagChecker } from '@/api/auth'
import { useCommandsApi } from '@/api/commands/commands.js'
import ActionConfirmation from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { toast } from 'vue-sonner'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{ row: Command }>()
const router = useRouter()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)

const manager = useCommandsApi()
const deleter = manager.useMutationDeleteCommand()
const patcher = manager.useMutationUpdateCommand()!

const { t } = useI18n()

const showDelete = ref(false)

async function switchEnabled(newValue: boolean) {
	await patcher?.executeMutation({
		id: props.row.id,
		opts: {
			enabled: newValue,
		},
	})

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})
}

async function deleteCommand() {
	await deleter.executeMutation({ id: props.row.id })

	toast.success(t('sharedTexts.deleted'), {
		duration: 1500,
	})
}

function goToCopyCommand() {
	router.push({
		path: '/dashboard/commands/custom/create',
		query: {
			commandIdForCopy: props.row.id,
		},
	})
}
</script>

<template>
	<div class="flex items-center gap-4">
		<Switch
			:disabled="!userCanManageCommands"
			:model-value="row.enabled"
			@update:model-value="switchEnabled"
		/>
		<div class="flex gap-2">
			<Tooltip v-if="row.module === 'CUSTOM'">
				<TooltipTrigger>
					<Button :disabled="!userCanManageCommands" size="icon" @click="goToCopyCommand">
						<CopyIcon class="h-4 w-4" />
					</Button>
				</TooltipTrigger>
				<TooltipContent>
					<p>Copy command as new</p>
				</TooltipContent>
			</Tooltip>

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

	<ActionConfirmation v-model:open="showDelete" @confirm="deleteCommand" />
</template>
