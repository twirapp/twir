import { createGlobalState } from '@vueuse/core'
import { ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import type { CommandsCreateOpts } from '@/gql/graphql'

import { useCommandsApi } from '@/api/commands/commands.js'
import { useToast } from '@/components/ui/toast'

type EditableCommand = CommandsCreateOpts & { id?: string, module: string }

const defaultFormValue: EditableCommand = {
	name: '',
	aliases: [],
	responses: [
		{
			text: '',
			order: 0,
		},
	],
	description: '',
	rolesIds: [],
	deniedUsersIds: [],
	allowedUsersIds: [],
	requiredMessages: 0,
	requiredUsedChannelPoints: 0,
	requiredWatchTime: 0,
	cooldown: 0,
	cooldownType: 'GLOBAL',
	isReply: true,
	visible: true,
	keepResponsesOrder: true,
	onlineOnly: false,
	enabled: true,
	groupId: null,
	cooldownRolesIds: [],
	enabledCategories: [],
	module: 'CUSTOM',
}

export const useCommandEdit = createGlobalState(() => {
	const commandsManager = useCommandsApi()
	const { toast } = useToast()
	const { t } = useI18n()

	const create = commandsManager.useMutationCreateCommand()
	const update = commandsManager.useMutationUpdateCommand()

	const { data: commands } = commandsManager.useQueryCommands()

	const formValue = ref<EditableCommand | null>(null)
	const isOpened = ref(false)

	function close() {
		isOpened.value = false
	}

	function editCommand(id: string) {
		const command = commands.value?.commands.find((command) => command.id === id)
		if (!command) {
			throw new Error(`Command with id ${id} not found`)
		}

		// for not modify original query object of command
		formValue.value = structuredClone(toRaw(command))
		isOpened.value = true
	}

	function createCommand() {
		formValue.value = structuredClone(defaultFormValue)
		isOpened.value = true
	}

	async function save() {
		if (!formValue.value) {
			throw new Error('Form value is not set')
		}

		const transformedOpts = {
			...formValue.value,
			// need to delete that
			default: undefined,
			defaultName: undefined,
			group: undefined,
			id: undefined,
			module: undefined,
			responses: formValue.value.responses.map((response, i) => ({
				text: response.text,
				order: i,
			})),
		}

		if (formValue.value.id) {
			await update.executeMutation({
				id: formValue.value.id,
				opts: transformedOpts,
			})
		} else {
			await create.executeMutation({
				opts: transformedOpts,
			})
		}

		close()

		toast({
			title: t('sharedTexts.saved'),
			variant: 'success',
			duration: 1500,
		})
	}

	return {
		isOpened,
		formValue,
		editCommand,
		createCommand,
		close,
		save,
	}
})
