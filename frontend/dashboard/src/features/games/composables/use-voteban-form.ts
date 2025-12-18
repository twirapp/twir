import { toTypedSchema } from '@vee-validate/zod';
import { createGlobalState } from '@vueuse/core';
import { useI18n } from 'vue-i18n';
import { z } from 'zod';

import { useGamesApi } from '@/api/games/games';
import { toast } from 'vue-sonner';
import { VoteBanGameVotingMode } from '@/gql/graphql';

const rules = z.object({
	enabled: z.boolean(),
	initMessage: z.string().max(500),
	banMessage: z.string().max(500),
	banMessageModerators: z.string().max(500),
	surviveMessage: z.string().max(500),
	surviveMessageModerators: z.string().max(500),
	votingMode: z.nativeEnum(VoteBanGameVotingMode),
	chatVotesWordsPositive: z.array(z.string()).max(10),
	chatVotesWordsNegative: z.array(z.string()).max(10),
	voteDuration: z.number().min(1).max(86400),
	neededVotes: z.number().min(1).max(999999),
	timeoutSeconds: z.number().min(1).max(86400),
	timeoutModerators: z.boolean(),
}).superRefine((data, ctx) => {
	const positives = new Set(data.chatVotesWordsPositive)

	for (const word of data.chatVotesWordsNegative) {
		if (positives.has(word)) {
			ctx.addIssue({
				code: 'custom',
				message: `Word "${word}" is already in use`,
				path: ['chatVotesWordsPositive'],
			})
			ctx.addIssue({
				code: 'custom',
				message: `Word "${word}" is already in use`,
				path: ['chatVotesWordsNegative'],
			})
			return
		}
	}
})

export const formSchema = toTypedSchema(rules)

export type FormSchema = z.infer<typeof rules>

export const useVotebanForm = createGlobalState(() => {
	const { t } = useI18n()
	const gamesApi = useGamesApi()
	const { data: settings } = gamesApi.useGamesQuery()
	const updater = gamesApi.useVotebanMutation()

	const initialValues: FormSchema = {
		enabled: false,
		initMessage:
			'The Twitch Police have decided that {targetUser} is not worthy of being in chat for not knowing memes. Write "{positiveTexts}" to support, or "{negativeTexts}" if you disagree.',
		banMessage: 'User {targetUser} is not worthy of being in chat.',
		banMessageModerators: 'User {targetUser} is not worthy of being in chat.',
		surviveMessage:
			'Looks like something is mixed up, {targetUser} is the kindest and most knowledgeable chat user.',
		surviveMessageModerators:
			'Looks like something is mixed up, {targetUser} is the kindest and most knowledgeable chat user.',
		votingMode: VoteBanGameVotingMode.Chat,
		chatVotesWordsPositive: ['Yay'],
		chatVotesWordsNegative: ['Nay'],
		voteDuration: 60,
		neededVotes: 3,
		timeoutSeconds: 60,
		timeoutModerators: false,
	}

	async function save(values: FormSchema) {
		const result = await updater.executeMutation({
			opts: values,
		})

		if (result.error) {
			toast.error(result.error.graphQLErrors?.map((e) => e.message).join(', ') ?? 'error', {
				duration: 5000,
			})
			return
		}

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
