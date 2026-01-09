import { createGlobalState } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { type TypeOf, array, boolean, nativeEnum, number, object, string } from 'zod'

import { useTimersApi } from '#layers/dashboard/api/timers'
import { TwitchAnnounceColor } from '~/gql/graphql.ts'

export const formSchema = object({
	id: string().optional(),
	name: string().min(2).max(50),
	timeInterval: number().int().min(0).max(1000).default(0),
	messageInterval: number().int().min(0).max(5000).default(0),
	responses: array(
		object({
			text: string().min(1).max(1000),
			isAnnounce: boolean(),
			count: number().int().min(1).max(20).default(1),
			announceColor: nativeEnum(TwitchAnnounceColor).default(TwitchAnnounceColor.Primary),
		})
	).min(1),
	enabled: boolean().default(true),
})

type FormSchema = TypeOf<typeof formSchema>

export const useTimersEdit = createGlobalState(() => {
	const { t } = useI18n()
	const router = useRouter()

	const timersApi = useTimersApi()
	const timers = timersApi.useQueryTimers()
	const updateMutation = timersApi.useMutationUpdateTimer()
	const createMutation = timersApi.useMutationCreateTimer()

	async function findTimer(id: string) {
		if (id === 'create') return

		const fetchedData = await timers.then((timers) => timers)
		const timer = fetchedData.data?.value?.timers.find((timer) => timer.id === id)

		if (!timer) throw new Error('Timer not found')

		return timer
	}

	async function submit(data: FormSchema) {
		if (data.id) {
			await updateMutation.executeMutation({
				id: data.id,
				opts: {
					...data,
					// this deletes id field from object, because backend will respond with error if it's presented
					//
					// @ts-expect-error
					id: undefined,
				},
			})
		} else {
			const result = await createMutation.executeMutation({
				opts: data,
			})

			if (result.error) {
				toast.error(result.error.graphQLErrors?.map((e) => e.message).join(', ') ?? 'error', {
					duration: 5000,
				})
				return
			}

			await router.push({
				path: `/dashboard/timers/${result.data?.timersCreate.id}`,
				state: {
					noTransition: true,
				},
			})
		}

		toast.success(t('sharedTexts.saved'), {
			duration: 2500,
		})
	}

	return {
		findTimer,
		timers,
		submit,
	}
})
