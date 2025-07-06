import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { KappagenEmojiStyle } from '@/gql/graphql'

export const kappagenFormSchema = toTypedSchema(
	z.object({
		enableSpawn: z.boolean().default(true),
		excludedEmotes: z.array(z.string()).default([]),
		enableRave: z.boolean().default(false),
		animation: z.object({
			fadeIn: z.boolean().default(true),
			fadeOut: z.boolean().default(true),
			zoomIn: z.boolean().default(false),
			zoomOut: z.boolean().default(false),
		}),
		animations: z
			.array(
				z.object({
					style: z.string().min(1, 'Animation style is required'),
					prefs: z.object({
						size: z.number().min(0.1).max(10).default(1),
						center: z.boolean().default(false),
						speed: z.number().min(1).max(100).default(50),
						faces: z.boolean().default(false),
						message: z.array(z.string()).default([]),
						time: z.number().min(100).max(30000).default(5000),
					}),
					count: z.number().min(1).max(100).default(1),
					enabled: z.boolean().default(true),
				})
			)
			.default([]),
		emotes: z.object({
			time: z.number().min(1000).max(60000).default(5000),
			max: z.number().min(1).max(500).default(100),
			queue: z.number().min(1).max(1000).default(100),
			ffzEnabled: z.boolean().default(true),
			bttvEnabled: z.boolean().default(true),
			sevenTvEnabled: z.boolean().default(true),
			emojiStyle: z.nativeEnum(KappagenEmojiStyle).default(KappagenEmojiStyle.Twemoji),
		}),
		size: z.object({
			rationNormal: z.number().min(0.1).max(5).default(1),
			rationSmall: z.number().min(0.1).max(5).default(0.5),
			min: z.number().min(10).max(200).default(20),
			max: z.number().min(50).max(500).default(150),
		}),
		events: z
			.array(
				z.object({
					event: z.string().min(1, 'Event type is required'),
					disabledAnimations: z.array(z.string()).default([]),
					enabled: z.boolean().default(true),
				})
			)
			.default([]),
	})
)
