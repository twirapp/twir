import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { EventOperationType, EventType } from '@/gql/graphql.ts'

export const eventFormSchema = toTypedSchema(z.object({
	type: z.nativeEnum(EventType),
	description: z.string().min(1).max(20),
	enabled: z.boolean().default(true),
	onlineOnly: z.boolean().default(false),
	rewardId: z.string().max(50).optional(),
	commandId: z.string().max(50),
	keywordId: z.string().max(50).optional(),
	operations: z.array(z.object({
		type: z.nativeEnum(EventOperationType),
		input: z.string().max(1000).optional(),
		delay: z.number().max(1000).min(0).default(0),
		repeat: z.number().min(0).max(50).default(0),
		useAnnounce: z.boolean().default(false),
		timeoutTime: z.number().min(0).default(0),
		timeoutMessage: z.string().optional(),
		target: z.string().max(1000).optional(),
		enabled: z.boolean().default(true),
		filters: z.array(z.object({
			type: z.string().max(1000).min(1),
			left: z.string().max(1000).min(1),
			right: z.string().max(1000).min(1),
		})).max(10).default([]),
	})).max(10).default([]),
}).superRefine((data, ctx) => {
	if (data.type === EventType.CommandUsed && !data.commandId) {
		ctx.addIssue({
			code: 'custom',
			message: 'Command ID is required',
			path: ['commandId'],
		})

		return false
	}
	if (data.type === EventType.KeywordMatched && !data.keywordId) {
		ctx.addIssue({
			code: 'custom',
			message: 'Keyword ID is required',
			path: ['keywordId'],
		})

		return false
	}

	if (data.type === EventType.KeywordMatched && data.keywordId) {
		ctx.addIssue({
			code: 'custom',
			message: 'Keyword ID is required',
			path: ['keywordId'],
		})

		return false
	}

	return true
}))
