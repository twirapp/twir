import { toTypedSchema } from '@vee-validate/zod'
import { createGlobalState } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { z } from 'zod'

import { useGamesApi } from '#layers/dashboard/api/games/games'

const rules = z.object({
	enabled: z.boolean(),
	startMessage: z.string().max(400),
	resultMessage: z.string().max(400),
	bothDieMessage: z.string().max(400),
	userCooldown: z.number().min(0).max(84000),
	globalCooldown: z.number().min(0).max(84000),
	secondsToAccept: z.number().min(0).max(3600),
	timeoutSeconds: z.number().min(0).max(84000),
	pointsPerWin: z.number().min(0).max(999999),
	pointsPerLose: z.number().min(0).max(999999),
	bothDiePercent: z.number().min(0).max(100),
})

export const formSchema = toTypedSchema(rules)

export type FormSchema = z.infer<typeof rules>

export const useDuelForm = createGlobalState(() => {
	const { t } = useI18n()
	const gamesApi = useGamesApi()
	const { data: settings } = gamesApi.useGamesQuery()
	const updater = gamesApi.useDuelMutation()

	const initialValues: FormSchema = {
		enabled: false,
		startMessage:
			'@{target}, @{initiator} challenges you to a fight. Use {duelAcceptCommandName} for next {acceptSeconds} seconds to accept the challenge.',
		resultMessage:
			"Sadly, @{loser} couldn't find a way to dodge the bullet and falls apart into eternal slumber.",
		bothDieMessage:
			'Unexpectedly @{initiator} and @{target} shoot each other. Only the time knows why this happened...',
		userCooldown: 0,
		globalCooldown: 0,
		secondsToAccept: 60,
		timeoutSeconds: 600,
		pointsPerWin: 0,
		pointsPerLose: 0,
		bothDiePercent: 0,
	}

	async function save(values: FormSchema) {
		await updater.executeMutation({
			opts: values,
		})

		toast.success(t('sharedTexts.saved'), {
			duration: 2500,
		})
	}

	return {
		initialValues,
		save,
		settings,
	}
})
