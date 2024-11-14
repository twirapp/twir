import { createGlobalState } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { array, boolean, nativeEnum, number, object, string } from 'zod'

import type { TypeOf } from 'zod'

import { useCommandsApi } from '@/api/commands/commands'
import { useRoles } from '@/api/roles'
import { useToast } from '@/components/ui/toast'
import { CommandExpiresType } from '@/gql/graphql'

export const formSchema = object({
	name: string().min(1).max(50),
	aliases: array(string().max(50)).max(50),
	responses: array(
		object({
			text: string().min(1).max(500),
			order: number(),
			twitchCategoriesIds: array(string()).max(100),
		}),
	).min(1),
	enabled: boolean().default(true),
	description: string().max(500),
	rolesIds: array(string()).max(100),
	deniedUsersIds: array(string()).max(100),
	allowedUsersIds: array(string()).max(100),
	requiredMessages: number().int().min(0).max(9999999999),
	requiredUsedChannelPoints: number().int().min(0).max(999999999999),
	requiredWatchTime: number().int().min(0).max(999999999999),
	cooldown: number().int().min(0).max(84600),
	cooldownType: string().optional(),
	cooldownRolesIds: array(string()).max(100),
	isReply: boolean().default(true),
	visible: boolean().default(true),
	keepResponsesOrder: boolean().default(true),
	onlineOnly: boolean().default(false),
	groupId: string().optional(),
	enabledCategories: array(string()).max(100),
	module: string().optional(),
	expiresAt: number().nullable().optional(),
	expiresType: nativeEnum(CommandExpiresType).nullable().optional(),
})

type FormSchema = TypeOf<typeof formSchema>

export const useCommandEditV2 = createGlobalState(() => {
	const { toast } = useToast()
	const { t } = useI18n()

	const commandsApi = useCommandsApi()
	const commands = commandsApi.useQueryCommands()
	const update = commandsApi.useMutationUpdateCommand()
	const create = commandsApi.useMutationCreateCommand()

	const rolesManager = useRoles()
	const { data: roles } = rolesManager.useRolesQuery()

	async function findCommand(id: string) {
		if (id === 'create') return

		const fetchedData = await commands.then((c) => c)
		const command = fetchedData.data?.value?.commands.find((command) => command.id === id)

		if (!command) throw new Error('Command not found')

		return command
	}

	async function submit(data: FormSchema) {
		if (data.id) {
			await update.executeMutation({
				id: data.id,
				opts: {
					...data,
					id: undefined,
				},
			})
		} else {
			await create.executeMutation({
				opts: data,
			})
		}

		toast({
			title: t('common.saved'),
			status: 'success',
		})
	}

	return {
		findCommand,
		submit,
		channelRoles: roles,
	}
})
