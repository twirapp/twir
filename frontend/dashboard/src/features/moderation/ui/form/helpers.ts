import {
	IconAbc,
	IconAsteriskSimple,
	IconLanguageOff,
	IconLinkOff,
	IconListLetters,
	IconMessageOff,
	IconMoodOff,
	type SVGProps,
} from '@tabler/icons-vue'
import { type FunctionalComponent, ref } from 'vue'

import type { Item as GrpcItem, ItemCreateMessage } from '@twir/api/messages/moderation/moderation'

export type Item = ItemCreateMessage & {
	icon: FunctionalComponent<SVGProps>
}

export type ItemData = Omit<GrpcItem, 'createdAt' | 'channelId' | 'updatedAt'>
export interface ItemWithOptionalId {
	id?: string
	data?: ItemData
}

export const Icons: Readonly<Record<string, (props: SVGProps) => FunctionalComponent<SVGProps, any, any>>> = Object.freeze({
	links: IconLinkOff,
	language: IconLanguageOff,
	deny_list: IconListLetters,
	long_message: IconMessageOff,
	caps: IconAbc,
	emotes: IconMoodOff,
	symbols: IconAsteriskSimple,
})

export const availableSettings: ItemData[] = [
	{
		deniedChatLanguages: [],
		banMessage: 'No links allowed',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'links',
		warningMessage: 'No links allowed [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Language not allowed',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'language',
		warningMessage: 'Language not allowed [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Bad word',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'deny_list',
		warningMessage: 'Bad word [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Too long message',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'long_message',
		warningMessage: 'Too long message [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Too much caps',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'caps',
		warningMessage: 'Too much caps [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Too many emotes',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 10,
		type: 'emotes',
		warningMessage: 'Too many emotes [warning]',
	},
	{
		deniedChatLanguages: [],
		banMessage: 'Too many symbols',
		banTime: 600,
		checkClips: false,
		denyList: [],
		enabled: true,
		excludedRoles: [],
		maxPercentage: 0,
		maxWarnings: 3,
		triggerLength: 0,
		type: 'symbols',
		warningMessage: 'Too many symbols [warning]',
	},
]

export const availableSettingsTypes = availableSettings.map(n => n.type)

const editableItem = ref<ItemWithOptionalId>()

export function useEditableItem() {
	const reset = () => editableItem.value = undefined

	return {
		editableItem,
		reset,
	}
}
