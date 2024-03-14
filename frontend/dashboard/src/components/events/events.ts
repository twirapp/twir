import {
	IconAccessPoint,
	IconAccessPointOff,
	IconAward,
	IconBracketsContain,
	IconCashBanknote,
	IconDeviceDesktopAnalytics,
	IconDice6,
	IconEraser,
	IconGift,
	IconHeartHandshake,
	IconMessageExclamation,
	IconPick,
	IconStar,
	IconTransform,
	IconUserCancel,
	IconUserDollar,
	IconUserHeart,
	IconUserPlus,
	IconUserStar,
} from '@tabler/icons-vue';
import { TwirEventType } from '@twir/api/messages/events/events';
import { FunctionalComponent } from 'vue';

export type TwirEvent = {
	name: string,
	icon?: FunctionalComponent,
	variables?: string[],
	type?: 'group',
	childrens?: Record<string, TwirEvent>,
	enumValue?: TwirEventType,
}

export const TWIR_EVENTS: Record<string, TwirEvent> = {
	FOLLOW: {
		name: 'Follow',
		icon: IconUserHeart,
		variables: ['userName', 'userDisplayName'],
		enumValue: TwirEventType.FOLLOW,
	},

	SUBS: {
		name: 'Subscribtions',
		type: 'group',
		childrens: {
			SUBSCRIBE: {
				name: 'Subscribe',
				icon: IconUserDollar,
				variables: ['userName', 'userDisplayName', 'subLevel'],
				enumValue: TwirEventType.SUBSCRIBE,
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
				enumValue: TwirEventType.RESUBSCRIBE,
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
				enumValue: TwirEventType.SUB_GIFT,
			},
		},
	},
	REDEMPTION_CREATED: {
		name: 'Reward Activated',
		icon: IconAward,
		variables: ['userName', 'userDisplayName', 'rewardName', 'rewardCost', 'rewardInput'],
		enumValue: TwirEventType.REDEMPTION_CREATED,
	},
	COMMAND_USED: {
		name: 'Command used',
		icon: IconStar,
		variables: ['userName', 'userDisplayName', 'commandName', 'commandInput'],
		enumValue: TwirEventType.COMMAND_USED,
	},
	FIRST_USER_MESSAGE: {
		name: 'First User Message',
		icon: IconMessageExclamation,
		variables: ['userName', 'userDisplayName'],
		enumValue: TwirEventType.FIRST_USER_MESSAGE,
	},

	STREAM: {
		name: 'Stream',
		type: 'group',
		childrens: {
			STREAM_FIRST_USER_JOIN: {
				name: 'User Joined Stream For First Time',
				icon: IconUserPlus,
				variables: ['userName'],
			},
			RAIDED: {
				name: 'Raided',
				icon: IconPick,
				variables: ['userName', 'userDisplayName', 'raidViewers'],
				enumValue: TwirEventType.RAIDED,
			},
			TITLE_OR_CATEGORY_CHANGED: {
				name: 'Title or Category Changed',
				icon: IconTransform,
				variables: [
					'oldStreamTitle',
					'newStreamTitle',
					'oldStreamCategory',
					'newStreamCategory',
				],
				enumValue: TwirEventType.TITLE_OR_CATEGORY_CHANGED,
			},
			STREAM_ONLINE: {
				name: 'Stream Online',
				icon: IconAccessPoint,
				variables: ['streamTitle', 'streamCategory'],
				enumValue: TwirEventType.STREAM_ONLINE,
			},
			STREAM_OFFLINE: {
				name: 'Stream Offline',
				icon: IconAccessPointOff,
				variables: [],
				enumValue: TwirEventType.STREAM_OFFLINE,
			},
		},
	},

	ON_CHAT_CLEAR: {
		name: 'On Chat Clear',
		icon: IconEraser,
		variables: [],
		enumValue: TwirEventType.CHAT_CLEAR,
	},
	DONATE: {
		name: 'Donate',
		icon: IconCashBanknote,
		variables: ['userName', 'donateAmount', 'donateCurrency', 'donateMessage'],
		enumValue: TwirEventType.DONATE,
	},
	KEYWORD_MATCHED: {
		name: 'Keyword Matched',
		icon: IconBracketsContain,
		variables: ['userName', 'userDisplayName', 'keywordName', 'keywordResponse'],
		enumValue: TwirEventType.KEYWORD_USED,
	},
	GREETING_SENDED: {
		name: 'Greeting Sended',
		icon: IconHeartHandshake,
		variables: ['userName', 'userDisplayName', 'greetingText'],
		enumValue: TwirEventType.GREETING_SENDED,
	},

	POLLS: {
		name: 'Polls',
		type: 'group',
		childrens: {
			POLL_BEGIN: {
				name: 'Poll Begin',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames'],
				enumValue: TwirEventType.POLL_STARTED,
			},
			POLL_PROGRESS: {
				name: 'Poll Progress',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames', 'pollTotalVotes'],
				enumValue: TwirEventType.POLL_VOTED,
			},
			POLL_END: {
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
				enumValue: TwirEventType.POLL_ENDED,
			},
		},
	},

	PREDICTIONS: {
		name: 'Predictions',
		type: 'group',
		childrens: {
			PREDICTION_BEGIN: {
				name: 'Prediction Begin',
				icon: IconDice6,
				variables: ['predictionTitle', 'predictionOptionsNames'],
				enumValue: TwirEventType.PREDICTION_STARTED,
			},
			PREDICTION_PROGRESS: {
				name: 'Prediction Progress',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
				enumValue: TwirEventType.PREDICTION_VOTED,
			},
			PREDICTION_LOCK: {
				name: 'Prediction Lock',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
				enumValue: TwirEventType.PREDICTION_LOCKED,
			},
			PREDICTION_END: {
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
				enumValue: TwirEventType.PREDICTION_ENDED,
			},
		},
	},

	CHANNEL_BAN: {
		name: 'Bans',
		type: 'group',
		childrens: {
			CHANNEL_BAN: {
				name: 'User banned/timeouted',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'moderatorName', 'moderatorDisplayName', 'banReason', 'banEndsInMinutes'],
				enumValue: TwirEventType.USER_BANNED,
			},
			CHANNEL_UNBAN_REQUEST_CREATE: {
				name: 'User Unban Request Created',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'message'],
			},
			CHANNEL_UNBAN_REQUEST_RESOLVE: {
				name: 'User Unban Request Accepted/Declined',
				icon: IconUserCancel,
				variables: ['userName', 'userDisplayName', 'moderatorName', 'moderatorDisplayName', 'message'],
			},
		},
	},
};
