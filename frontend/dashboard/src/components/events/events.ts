import { IconAccessPoint, IconAccessPointOff, IconAward, IconBracketsContain, IconCashBanknote, IconDeviceDesktopAnalytics, IconDice6, IconEraser, IconGift, IconHeartHandshake, IconMessageExclamation, IconPick, IconStar, IconTransform, IconUserCancel, IconUserDollar, IconUserHeart, IconUserStar } from '@tabler/icons-vue'

import type { FunctionalComponent } from 'vue'

import { EventType } from '@/api/events.ts'

export interface TwirEvent {
	name: string
	icon?: FunctionalComponent
	variables?: string[]
	type?: 'group'
	childrens?: Record<string, TwirEvent>
	enumValue?: EventType
}

export const EventsOptions: Record<string, TwirEvent> = {
	[EventType.Follow]: {
		name: 'Follow',
		icon: IconUserHeart,
		variables: ['userName', 'userDisplayName'],
		enumValue: EventType.Follow,
	},

	SUBS: {
		name: 'Subscribtions',
		type: 'group',
		childrens: {
			SUBSCRIBE: {
				name: 'Subscribe',
				icon: IconUserDollar,
				variables: ['userName', 'userDisplayName', 'subLevel'],
				enumValue: EventType.Subscribe,
			},
			RESUBSCRIBE: {
				name: 'Resubscribe',
				icon: IconUserStar,
				variables: [
					'userName',
					'userDisplayName',
					'subLevel',
					'resubMonths',
					'resubStreak',
					'resubMessage',
				],
				enumValue: EventType.Resubscribe,
			},
			SUB_GIFT: {
				name: 'Subscribe Gift',
				icon: IconGift,
				variables: [
					'userName',
					'userDisplayName',
					'targetUserName',
					'targetDisplayName',
					'subLevel',
				],
				enumValue: EventType.SubGift,
			},
		},
	},
	[EventType.RedemptionCreated]: {
		name: 'Reward Activated',
		icon: IconAward,
		variables: ['userName', 'userDisplayName', 'rewardName', 'rewardCost', 'rewardInput'],
		enumValue: EventType.RedemptionCreated,
	},
	[EventType.CommandUsed]: {
		name: 'Command used',
		icon: IconStar,
		variables: ['userName', 'userDisplayName', 'commandName', 'commandInput'],
		enumValue: EventType.CommandUsed,
	},
	[EventType.FirstUserMessage]: {
		name: 'First User Message',
		icon: IconMessageExclamation,
		variables: ['userName', 'userDisplayName'],
		enumValue: EventType.FirstUserMessage,
	},

	STREAM: {
		name: 'Stream',
		type: 'group',
		childrens: {
			// [EventType.FirstUserMessage]: {
			// 	name: 'User Joined Stream For First Time',
			// 	icon: IconUserPlus,
			// 	variables: ['userName'],
			// 	enumValue: EventType.FirstUserMessage,
			// },
			[EventType.Raided]: {
				name: 'Raided',
				icon: IconPick,
				variables: ['userName', 'userDisplayName', 'raidViewers'],
				enumValue: EventType.Raided,
			},
			[EventType.TitleOrCategoryChanged]: {
				name: 'Title or Category Changed',
				icon: IconTransform,
				variables: [
					'oldStreamTitle',
					'newStreamTitle',
					'oldStreamCategory',
					'newStreamCategory',
				],
				enumValue: EventType.TitleOrCategoryChanged,
			},
			[EventType.StreamOnline]: {
				name: 'Stream Online',
				icon: IconAccessPoint,
				variables: ['streamTitle', 'streamCategory'],
				enumValue: EventType.StreamOnline,
			},
			[EventType.StreamOffline]: {
				name: 'Stream Offline',
				icon: IconAccessPointOff,
				variables: [],
				enumValue: EventType.StreamOffline,
			},
		},
	},

	[EventType.OnChatClear]: {
		name: 'On Chat Clear',
		icon: IconEraser,
		variables: [],
		enumValue: EventType.OnChatClear,
	},
	[EventType.Donate]: {
		name: 'Donate',
		icon: IconCashBanknote,
		variables: ['userName', 'donateAmount', 'donateCurrency', 'donateMessage'],
		enumValue: EventType.Donate,
	},
	[EventType.KeywordMatched]: {
		name: 'Keyword Matched',
		icon: IconBracketsContain,
		variables: ['userName', 'userDisplayName', 'keywordName', 'keywordResponse'],
		enumValue: EventType.KeywordMatched,
	},
	[EventType.GreetingSended]: {
		name: 'Greeting Sended',
		icon: IconHeartHandshake,
		variables: ['userName', 'userDisplayName', 'greetingText'],
		enumValue: EventType.GreetingSended,
	},

	POLLS: {
		name: 'Polls',
		type: 'group',
		childrens: {
			[EventType.PollBegin]: {
				name: 'Poll Begin',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames'],
				enumValue: EventType.PollBegin,
			},
			[EventType.PollProgress]: {
				name: 'Poll Progress',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames', 'pollTotalVotes'],
				enumValue: EventType.PollProgress,
			},
			[EventType.PollEnd]: {
				name: 'Poll End',
				icon: IconDeviceDesktopAnalytics,
				variables: [
					'pollTitle',
					'pollOptionsNames',
					'pollTotalVotes',
					'pollWinnerTitle',
					'pollWinnerBitsVotes',
					'pollWinnerChannelsPointsVotes',
					'pollWinnerTotalVotes',
				],
				enumValue: EventType.PollEnd,
			},
		},
	},

	PREDICTIONS: {
		name: 'Predictions',
		type: 'group',
		childrens: {
			[EventType.PredictionBegin]: {
				name: 'Prediction Begin',
				icon: IconDice6,
				variables: ['predictionTitle', 'predictionOptionsNames'],
				enumValue: EventType.PredictionBegin,
			},
			[EventType.PredictionProgress]: {
				name: 'Prediction Progress',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
				enumValue: EventType.PredictionProgress,
			},
			[EventType.PredictionLock]: {
				name: 'Prediction Lock',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
				enumValue: EventType.PredictionLock,
			},
			[EventType.PredictionEnd]: {
				name: 'Prediction End',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
					`predictionWinner.title`,
					`predictionWinner.totalUsers`,
					`predictionWinner.totalPoints`,
					`predictionWinner.topUsers`,
				],
				enumValue: EventType.PredictionEnd,
			},
		},
	},

	CHANNEL_BAN: {
		name: 'Bans',
		type: 'group',
		childrens: {
			[EventType.ChannelBan]: {
				name: 'User banned/timeouted',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'moderatorName', 'moderatorDisplayName', 'banReason', 'banEndsInMinutes'],
				enumValue: EventType.ChannelBan,
			},
			[EventType.ChannelUnbanRequestCreate]: {
				name: 'User Unban Request Created',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'message'],
				enumValue: EventType.ChannelUnbanRequestCreate,
			},
			[EventType.ChannelUnbanRequestResolve]: {
				name: 'User Unban Request Accepted/Declined',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'moderatorName', 'moderatorDisplayName', 'message'],
				enumValue: EventType.ChannelUnbanRequestResolve,
			},
		},
	},
}
