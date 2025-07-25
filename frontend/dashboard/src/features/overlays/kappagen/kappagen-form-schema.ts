import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { KappagenEmojiStyle, KappagenOverlayAnimationStyle } from '@/gql/graphql'

const schema = z.object({
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
				style: z.nativeEnum(KappagenOverlayAnimationStyle),
				prefs: z
					.object({
						size: z.number().min(0.1).max(10).default(1).nullable(),
						center: z.boolean().default(false).nullable(),
						speed: z.number().min(1).max(100).default(50).nullable(),
						faces: z.boolean().default(false).nullable(),
						message: z.array(z.string()).default([]).nullable(),
						time: z.number().min(100).max(30000).default(5000).nullable(),
					})
					.nullable()
					.default(null),
				count: z.number().min(1).max(1000).nullable(),
				enabled: z.boolean().default(true),
			}),
		)
		.default([]),
	emotes: z.object({
		time: z.number().min(1).max(15).default(1),
		max: z.number().min(1).max(500).default(100),
		queue: z.number().min(1).max(1000).default(100),
		ffzEnabled: z.boolean().default(true),
		bttvEnabled: z.boolean().default(true),
		sevenTvEnabled: z.boolean().default(true),
		emojiStyle: z.nativeEnum(KappagenEmojiStyle).default(KappagenEmojiStyle.Twemoji),
	}),
	size: z.object({
		rationNormal: z.number().min(0.05).max(0.15).default(0.05),
		rationSmall: z.number().min(0.02).max(0.07).default(0.02),
		min: z.number().default(20),
		max: z.number().default(150),
	}),
	events: z
		.array(
			z.object({
				event: z.string().min(1, 'Event type is required'),
				disabledAnimations: z.array(z.string()).default([]),
				enabled: z.boolean().default(true),
			}),
		)
		.default([]),
})

export const kappagenFormSchema = toTypedSchema(schema)

export type KappagenFormSchema = z.infer<typeof schema>
