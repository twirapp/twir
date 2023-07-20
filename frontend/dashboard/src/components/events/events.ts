import {
	IconDeviceDesktopAnalytics,
	IconMessageExclamation,
	IconBracketsContain,
	IconAccessPointOff,
	IconHeartHandshake,
	IconCashBanknote,
	IconAccessPoint,
	IconUserDollar,
	IconTransform,
	IconUserHeart,
	IconUserPlus,
	IconUserStar,
	IconEraser,
	IconAward,
	IconDice6,
	IconStar,
	IconPick,
	IconGift,
} from '@tabler/icons-vue';
import { FunctionalComponent } from 'vue';

type TwirEvent = {
	name: string,
	icon?: FunctionalComponent,
	variables?: string[],
	type?: 'group',
	childrens?: Record<string, TwirEvent>
}

export const EVENTS: Record<string, TwirEvent> = {
	FOLLOW: {
		name: 'Follow',
		icon: IconUserHeart,
		variables: ['userName', 'userDisplayName'],
	},

	SUBS: {
		name: 'Subscribtions',
		type: 'group',
		childrens: {
			SUBSCRIBE: {
				name: 'Subscribe',
				icon: IconUserDollar,
				variables: ['userName', 'userDisplayName', 'subLevel'],
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
			},
		},
	},
	REDEMPTION_CREATED: {
		name: 'Reward Activated',
		icon: IconAward,
		variables: ['userName', 'userDisplayName', 'rewardName', 'rewardCost', 'rewardInput'],
	},
	COMMAND_USED: {
		name: 'Reward Activated',
		icon: IconStar,
		variables: ['userName', 'userDisplayName', 'commandName', 'commandInput'],
	},
	FIRST_USER_MESSAGE: {
		name: 'First User Message',
		icon: IconMessageExclamation,
		variables: ['userName', 'userDisplayName'],
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
			},
			STREAM_ONLINE: {
				name: 'Stream Online',
				icon: IconAccessPoint,
				variables: ['streamTitle', 'streamCategory'],
			},
			STREAM_OFFLINE: {
				name: 'Stream Offline',
				icon: IconAccessPointOff,
				variables: [],
			},
		},
	},
	
	ON_CHAT_CLEAR: {
		name: 'On Chat Clear',
		icon: IconEraser,
		variables: [],
	},
	DONATE: {
		name: 'Donate',
		icon: IconCashBanknote,
		variables: ['userName', 'donateAmount', 'donateCurrency', 'donateMessage'],
	},
	KEYWORD_MATCHED: {
		name: 'Keyword Matched',
		icon: IconBracketsContain,
		variables: ['userName', 'userDisplayName', 'keywordName', 'keywordResponse'],
	},
	GREETING_SENDED: {
		name: 'Greeting Sended',
		icon: IconHeartHandshake,
		variables: ['userName', 'userDisplayName', 'greetingText'],
	},

	POLLS: {
		name: 'Polls',
		type: 'group',
		childrens: {
			POLL_BEGIN: {
				name: 'Poll Begin',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames'],
			},
			POLL_PROGRESS: {
				name: 'Poll Progress',
				icon: IconDeviceDesktopAnalytics,
				variables: ['pollTitle', 'pollOptionsNames', 'pollTotalVotes'],
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
			},
			PREDICTION_PROGRESS: {
				name: 'Prediction Progress',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
			},
			PREDICTION_LOCK: {
				name: 'Prediction Lock',
				icon: IconDice6,
				variables: [
					'predictionTitle',
					'predictionOptionsNames',
					'predictionTotalChannelPoints',
				],
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
			},
		},
	},
};
