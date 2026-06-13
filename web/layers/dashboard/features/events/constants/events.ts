import {
	BanIcon,
	BanknoteArrowDownIcon,
	BubblesIcon,
	CommandIcon,
	DollarSignIcon,
	GiftIcon,
	HeartHandshakeIcon,
	HeartIcon,
	MessageCircleWarningIcon,
	PaintbrushIcon,
	PickaxeIcon,
	ReplaceAllIcon,
	StarIcon,
	TrendingUpDown,
	VideoIcon,
	VideoOffIcon,
	VoteIcon,
	WholeWordIcon,


import { EventType } from '~~/layers/dashboard/api/events.js'

export interface TwirEvent {
	name: string
	icon?: string
	variables?: string[]
	type?: 'group'
	childrens?: Record<string, TwirEvent>
	enumValue?: EventType
}

export const EventsOptions: Record<string, TwirEvent> = {
	[EventType.Follow]: {
		name: 'Follow',
		icon: 'lucide:heart',
		variables: ['userName', 'userDisplayName', 'channelFollowers', 'channelStreamFollowers'],
		enumValue: EventType.Follow,
	},

	SUBS: {
		name: 'Subscriptions',
		type: 'group',
		childrens: {
			SUBSCRIBE: {
				name: 'Subscribe',
				icon: 'lucide:dollar-sign',
				variables: ['userName', 'userDisplayName', 'subLevel'],
				enumValue: EventType.Subscribe,
			},
			RESUBSCRIBE: {
				name: 'Resubscribe',
				icon: 'lucide:star',
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
				icon: 'lucide:gift',
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
		icon: 'lucide:bubbles',
		variables: ['userName', 'userDisplayName', 'rewardName', 'rewardCost', 'rewardInput'],
		enumValue: EventType.RedemptionCreated,
	},
	[EventType.CommandUsed]: {
		name: 'Command used',
		icon: 'lucide:command',
		variables: ['userName', 'userDisplayName', 'commandName', 'commandInput'],
		enumValue: EventType.CommandUsed,
	},
	[EventType.FirstUserMessage]: {
		name: 'First User Message',
		icon: 'lucide:message-circle-warning',
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
				icon: 'lucide:pickaxe',
				variables: ['userName', 'userDisplayName', 'raidViewers'],
				enumValue: EventType.Raided,
			},
			[EventType.TitleOrCategoryChanged]: {
				name: 'Title or Category Changed',
				icon: 'lucide:replace-all',
				variables: ['oldStreamTitle', 'newStreamTitle', 'oldStreamCategory', 'newStreamCategory'],
				enumValue: EventType.TitleOrCategoryChanged,
			},
			[EventType.StreamOnline]: {
				name: 'Stream Online',
				icon: 'lucide:video',
				variables: ['streamTitle', 'streamCategory'],
				enumValue: EventType.StreamOnline,
			},
			[EventType.StreamOffline]: {
				name: 'Stream Offline',
				icon: 'lucide:video-off',
				variables: [],
				enumValue: EventType.StreamOffline,
			},
		},
	},

	[EventType.OnChatClear]: {
		name: 'On Chat Clear',
		icon: 'lucide:paintbrush',
		variables: [],
		enumValue: EventType.OnChatClear,
	},
	[EventType.Donate]: {
		name: 'Donate',
		icon: 'lucide:banknote-arrow-down',
		variables: ['userName', 'donateAmount', 'donateCurrency', 'donateMessage'],
		enumValue: EventType.Donate,
	},
	[EventType.KeywordMatched]: {
		name: 'Keyword Matched',
		icon: 'lucide:whole-word',
		variables: ['userName', 'userDisplayName', 'keywordName', 'keywordResponse'],
		enumValue: EventType.KeywordMatched,
	},
	[EventType.GreetingSended]: {
		name: 'Greeting Sended',
		icon: 'lucide:heart-handshake',
		variables: ['userName', 'userDisplayName', 'greetingText'],
		enumValue: EventType.GreetingSended,
	},

	POLLS: {
		name: 'Polls',
		type: 'group',
		childrens: {
			[EventType.PollBegin]: {
				name: 'Poll Begin',
				icon: 'lucide:vote',
				variables: ['pollTitle', 'pollOptionsNames'],
				enumValue: EventType.PollBegin,
			},
			[EventType.PollProgress]: {
				name: 'Poll Progress',
				icon: 'lucide:vote',
				variables: ['pollTitle', 'pollOptionsNames', 'pollTotalVotes'],
				enumValue: EventType.PollProgress,
			},
			[EventType.PollEnd]: {
				name: 'Poll End',
				icon: 'lucide:vote',
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
				icon: 'lucide:trending-up-down',
				variables: ['predictionTitle', 'predictionOptionsNames'],
				enumValue: EventType.PredictionBegin,
			},
			[EventType.PredictionProgress]: {
				name: 'Prediction Progress',
				icon: 'lucide:trending-up-down',
				variables: ['predictionTitle', 'predictionOptionsNames', 'predictionTotalChannelPoints'],
				enumValue: EventType.PredictionProgress,
			},
			[EventType.PredictionLock]: {
				name: 'Prediction Lock',
				icon: 'lucide:trending-up-down',
				variables: ['predictionTitle', 'predictionOptionsNames', 'predictionTotalChannelPoints'],
				enumValue: EventType.PredictionLock,
			},
			[EventType.PredictionEnd]: {
				name: 'Prediction End',
				icon: 'lucide:trending-up-down',
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
				icon: 'lucide:ban',
				variables: [
					'userName',
					'userDisplayName',
					'moderatorName',
					'moderatorDisplayName',
					'banReason',
					'banEndsInMinutes',
				],
				enumValue: EventType.ChannelBan,
			},
			[EventType.ChannelUnbanRequestCreate]: {
				name: 'User Unban Request Created',
				icon: 'lucide:ban',
				variables: ['userName', 'userDisplayName', 'message'],
				enumValue: EventType.ChannelUnbanRequestCreate,
			},
			[EventType.ChannelUnbanRequestResolve]: {
				name: 'User Unban Request Accepted/Declined',
				icon: 'lucide:ban',
				variables: [
					'userName',
					'userDisplayName',
					'moderatorName',
					'moderatorDisplayName',
					'message',
				],
				enumValue: EventType.ChannelUnbanRequestResolve,
			},
		},
	},
}
