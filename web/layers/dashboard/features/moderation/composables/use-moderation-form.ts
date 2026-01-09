import type { FunctionalComponent } from 'vue'

import { toTypedSchema } from '@vee-validate/zod'
import {
	AsteriskIcon,
	CaseUpperIcon,
	LanguagesIcon,
	LetterTextIcon,
	Link2OffIcon,
	Repeat1Icon,
	RulerDimensionLineIcon,
	SmileIcon,
} from 'lucide-vue-next'
import { z } from 'zod'

import type { ModerationItem } from '#layers/dashboard/api/moderation'

import { ModerationSettingsType } from '~/gql/graphql.ts'

export type EditableItem = Omit<
	Omit<ModerationItem, 'id'> & {
		id?: string
	},
	'createdAt' | 'channelId' | 'updatedAt'
>

export const Icons: Readonly<Record<ModerationSettingsType, FunctionalComponent>> = Object.freeze({
	[ModerationSettingsType.Links]: Link2OffIcon,
	[ModerationSettingsType.Language]: LanguagesIcon,
	[ModerationSettingsType.DenyList]: LetterTextIcon,
	[ModerationSettingsType.LongMessage]: RulerDimensionLineIcon,
	[ModerationSettingsType.Caps]: CaseUpperIcon,
	[ModerationSettingsType.Emotes]: SmileIcon,
	[ModerationSettingsType.Symbols]: AsteriskIcon,
	[ModerationSettingsType.OneManSpam]: Repeat1Icon,
})

export const moderationValidationRules = toTypedSchema(
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
)
