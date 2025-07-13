import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { EventOperationType, EventType } from '@/gql/graphql.ts'

export const eventFormSchema = toTypedSchema(
	z
		.object({
			type: z.nativeEnum(EventType),
			description: z.string().min(1).max(20),
			enabled: z.boolean().default(true),
			onlineOnly: z.boolean().default(false),
			rewardId: z.string().max(50).optional(),
			commandId: z.string().max(50).optional(),
			keywordId: z.string().max(50).optional(),
			operations: z
				.array(
					z.object({
						type: z.nativeEnum(EventOperationType),
						input: z.string().max(1000).optional(),
						delay: z.number().max(1000).min(0).default(0),
						repeat: z.number().min(0).max(50).default(0),
						useAnnounce: z.boolean().default(false),
						timeoutTime: z.number().min(0).default(0),
						timeoutMessage: z.string().optional(),
						target: z.string().max(1000).optional(),
						enabled: z.boolean().default(true),
						filters: z
							.array(
								z.object({
									type: z.string().max(1000).min(1),
									left: z.string().max(1000).min(1),
									right: z.string().max(1000).min(1),
								})
							)
							.max(10)
							.default([]),
					})
				)
				.max(10)
				.default([]),
		})
		.superRefine((data, ctx) => {
			let success = true

			if (data.type === EventType.CommandUsed && !data.commandId) {
				ctx.addIssue({
					code: 'custom',
					message: 'Command ID is required',
					path: ['commandId'],
				})

				success = false
			}
			if (data.type === EventType.KeywordMatched && !data.keywordId) {
				ctx.addIssue({
					code: 'custom',
					message: 'Keyword ID is required',
					path: ['keywordId'],
				})

				success = false
			}

			if (data.type === EventType.KeywordMatched && data.keywordId) {
				ctx.addIssue({
					code: 'custom',
					message: 'Keyword ID is required',
					path: ['keywordId'],
				})

				success = false
			}

			for (const operation of data.operations) {
				const index = data.operations.indexOf(operation)

				if (
					[
						EventOperationType.ChangeVariable,
						EventOperationType.DecrementVariable,
						EventOperationType.IncrementVariable,
					].includes(operation.type)
				) {
					if (!operation.target) {
						ctx.addIssue({
							code: 'custom',
							message: 'Target is required for variable operations',
							path: [`operations[${index}].target`],
						})
						success = false
					}

					if (!operation.input) {
						ctx.addIssue({
							code: 'custom',
							message: 'Input is required for variable operations',
							path: [`operations[${index}].input`],
						})
						success = false
					}
				}

				if (
					[
						EventOperationType.AllowCommandToUser,
						EventOperationType.RemoveAllowCommandToUser,
						EventOperationType.DenyCommandToUser,
						EventOperationType.RemoveDenyCommandToUser,
					].includes(operation.type)
				) {
					if (!operation.target) {
						ctx.addIssue({
							code: 'custom',
							message: 'Target is required for command operations',
							path: [`operations[${index}].target`],
						})
						success = false
					}

					if (!operation.input) {
						ctx.addIssue({
							code: 'custom',
							message: 'Input is required for command operations',
							path: [`operations[${index}].input`],
						})
						success = false
					}
				}
			}

			return success
		})
)
