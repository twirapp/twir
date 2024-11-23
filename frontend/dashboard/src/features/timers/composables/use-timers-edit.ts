import { createGlobalState } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { type TypeOf, array, boolean, number, object, string } from 'zod'

import { useTimersApi } from '@/api/timers'
import { useToast } from '@/components/ui/toast'

export const formSchema = object({
	id: string().optional(),
	name: string().min(2).max(50),
	timeInterval: number().int().min(0).max(100).default(0),
	messageInterval: number().int().min(0).max(5000).default(0),
	responses: array(
		object({
			text: string().min(1).max(500),
			isAnnounce: boolean(),
		}),
	).min(1),
	enabled: boolean().default(true),
})

type FormSchema = TypeOf<typeof formSchema>

export const useTimersEdit = createGlobalState(() => {
	const { toast } = useToast()
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
					// eslint-disable-next-line ts/ban-ts-comment
					// @ts-expect-error
					id: undefined,
				},
			})
		} else {
			const result = await createMutation.executeMutation({
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

			await router.push(`/dashboard/timers/${result.data?.timersCreate.id}`)
		}

		toast({
			title: t('sharedTexts.saved'),
			duration: 2500,
			variant: 'success',
		})
	}

	return {
		findTimer,
		timers,
		submit,
	}
})
