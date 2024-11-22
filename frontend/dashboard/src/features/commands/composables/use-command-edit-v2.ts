import { createGlobalState } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { array, boolean, nativeEnum, number, object, string } from 'zod'

import type { TypeOf } from 'zod'

import { type Command, useCommandsApi } from '@/api/commands/commands'
import { useRoles } from '@/api/roles'
import { useToast } from '@/components/ui/toast'
import { CommandExpiresType } from '@/gql/graphql'

export const formSchema = object({
	id: string().optional(),
	name: string().min(1).max(50),
	aliases: array(string().max(50)).max(50),
	enabled: boolean(),
	responses: array(
		object({
			text: string().min(1).max(500),
			twitchCategoriesIds: array(string()).max(100),
		}),
	).max(3).default([]),
	description: string().max(500),
	rolesIds: array(string()).max(100),
	deniedUsersIds: array(string()).max(100),
	allowedUsersIds: array(string()).max(100),
	requiredMessages: number().int().min(0).max(9999999999),
	requiredUsedChannelPoints: number().int().min(0).max(999999999999),
	requiredWatchTime: number().int().min(0).max(999999999999),
	cooldown: number().int().min(0).max(84600),
	cooldownType: string(),
	cooldownRolesIds: array(string()).max(100),
	isReply: boolean(),
	visible: boolean(),
	keepResponsesOrder: boolean(),
	onlineOnly: boolean(),
	groupId: string().nullable().optional().default(null),
	enabledCategories: array(string()).max(100),
	module: string().optional(),
})
	.and(object({
		expiresAt: number().nullable().optional(),
		expiresType: nativeEnum(CommandExpiresType).nullable().optional(),
	}).refine((data) => {
		if (data.expiresAt && !data.expiresType) {
			return false
		}

		if (!data.expiresAt && data.expiresType) {
			return false
		}

		return true
	}, {
		message: 'ExpiresAt and ExpiresType must be both set or both not set',
		path: ['expiresAt', 'expiresType'],
	}))

type FormSchema = TypeOf<typeof formSchema>

export const useCommandEditV2 = createGlobalState(() => {
	const { toast } = useToast()
	const { t } = useI18n()
	const router = useRouter()

	const commandsApi = useCommandsApi()
	const commands = commandsApi.useQueryCommands()
	const update = commandsApi.useMutationUpdateCommand()
	const create = commandsApi.useMutationCreateCommand()

	const rolesManager = useRoles()
	const { data: roles } = rolesManager.useRolesQuery()

	const command = ref<Command | null>(null)
	const isCustom = computed(() => {
		return !command.value?.default
	})

	async function findCommand(id: string) {
		command.value = null
		if (id === 'create') return

		const fetchedData = await commands.then((c) => c)
		const foundCmd = fetchedData.data?.value?.commands.find((command) => command.id === id)

		if (!foundCmd) {
			throw new Error('Command not found')
		}

		command.value = foundCmd

		return foundCmd
	}

	async function submit(data: FormSchema) {
		if (data.id) {
			await update.executeMutation({
				id: data.id,
				opts: {
					...data,
					// eslint-disable-next-line ts/ban-ts-comment
					// @ts-expect-error
					id: undefined,
					module: undefined,
				},
			})
		} else {
			const result = await create.executeMutation({
				opts: data,
			})

			if (result.error) {
				toast({
					title: result.error.graphQLErrors?.map(e => e.message).join(', ') ?? 'error',
					duration: 5000,
					variant: 'destructive',
				})
				return
			}

			await router.push(`/dashboard/commands/custom/${result.data?.commandsCreate.id}`)
		}

		toast({
			title: t('sharedTexts.saved'),
			duration: 2500,
			variant: 'success',
		})
	}

	return {
		findCommand,
		submit,
		channelRoles: roles,
		command,
		isCustom,
	}
})
