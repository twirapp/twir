
import { z } from 'zod'

import type { ModerationItem } from '~~/layers/dashboard/api/moderation'

import { ModerationSettingsType } from '~/gql/graphql.js'

export type EditableItem = Omit<
	Omit<ModerationItem, 'id'> & {
		id?: string
	},
	'createdAt' | 'channelId' | 'updatedAt'
>

export const Icons: Readonly<Record<ModerationSettingsType, string>> = Object.freeze({
	[ModerationSettingsType.Links]: 'lucide:link-2-off',
	[ModerationSettingsType.Language]: 'lucide:languages',
	[ModerationSettingsType.DenyList]: 'lucide:letter-text',
	[ModerationSettingsType.LongMessage]: 'lucide:ruler-dimension-line',
	[ModerationSettingsType.Caps]: 'lucide:case-upper',
	[ModerationSettingsType.Emotes]: 'lucide:smile',
	[ModerationSettingsType.Symbols]: 'lucide:asterisk',
	[ModerationSettingsType.OneManSpam]: 'lucide:repeat-1',
})

export const moderationValidationRules =
	z.object({
		id: z.string().optional(),
		name: z.string().optional().nullable(),
		banMessage: z.string().max(500).optional().default(''),
		banTime: z.number().min(0).max(86400).default(600),
		warningMessage: z.string().max(500).optional().default(''),
		maxWarnings: z.number().min(0).max(10).default(3),
		excludedRoles: z.array(z.string()).default([]),
		enabled: z.boolean().default(true),
		checkClips: z.boolean().default(true),
		triggerLength: z.number().min(0).max(10000).default(30),
		maxPercentage: z.number().min(0).max(100).default(50),
		denyList: z.array(z.string()).default([]),
		deniedChatLanguages: z.array(z.string()).default([]),
		type: z.nativeEnum(ModerationSettingsType),
		denyListRegexpEnabled: z.boolean().default(false),
		denyListWordBoundaryEnabled: z.boolean().default(false),
		denyListSensitivityEnabled: z.boolean().default(false),
		oneManSpamMinimumStoredMessages: z.number().min(0).max(20).default(5),
		oneManSpamMessageMemorySeconds: z.number().min(0).max(600).default(30),
		languageExcludedWords: z.array(z.string()).default([]),
	})
