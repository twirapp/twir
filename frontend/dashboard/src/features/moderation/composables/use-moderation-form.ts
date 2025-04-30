import {
	IconAbc,
	IconAsteriskSimple,
	IconLanguageOff,
	IconLinkOff,
	IconListLetters,
	IconMessageOff,
	IconMoodOff,
} from '@tabler/icons-vue'
import { toTypedSchema } from '@vee-validate/zod'
import { z } from 'zod'

import type { ModerationItem } from '@/api'
import type {
	SVGProps,
} from '@tabler/icons-vue'
import type { FunctionalComponent } from 'vue'

import { ModerationSettingsType } from '@/gql/graphql.ts'

export type EditableItem = Omit<Omit<ModerationItem, 'id'> & {
	id?: string
}, 'createdAt' | 'channelId' | 'updatedAt'>

export const Icons: Readonly<Record<ModerationSettingsType, (props: SVGProps) => FunctionalComponent<SVGProps, any, any>>> = Object.freeze({
	[ModerationSettingsType.Links]: IconLinkOff,
	[ModerationSettingsType.Language]: IconLanguageOff,
	[ModerationSettingsType.DenyList]: IconListLetters,
	[ModerationSettingsType.LongMessage]: IconMessageOff,
	[ModerationSettingsType.Caps]: IconAbc,
	[ModerationSettingsType.Emotes]: IconMoodOff,
	[ModerationSettingsType.Symbols]: IconAsteriskSimple,
})

export const moderationValidationRules = toTypedSchema(z.object({
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
}))
