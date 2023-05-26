import { EventType } from '@tsuwari/typeorm/entities/events/Event';

export const eventsMapping: Record<
	keyof typeof EventType,
	{
		description: string;
		availableVariables: Array<string>;
	}
> = {
	FOLLOW: {
		description: '',
		availableVariables: ['userName', 'userDisplayName'],
	},
	SUBSCRIBE: {
		description: '',
		availableVariables: ['userName', 'userDisplayName', 'subLevel'],
	},
	RESUBSCRIBE: {
		description: '',
		availableVariables: [
			'userName',
			'userDisplayName',
			'subLevel',
			'resubMonths',
			'resubStreak',
			'resubMessage',
		],
	},
	SUB_GIFT: {
		description: '',
		availableVariables: [
			'userName',
			'userDisplayName',
			'targetUserName',
			'targetDisplayName',
			'subLevel',
		],
	},
	REDEMPTION_CREATED: {
		description: 'Reward activated',
		availableVariables: ['userName', 'userDisplayName', 'rewardName', 'rewardCost', 'rewardInput'],
	},
	COMMAND_USED: {
		description: '',
		availableVariables: ['userName', 'userDisplayName', 'commandName', 'commandInput'],
	},
	FIRST_USER_MESSAGE: {
		description: '',
		availableVariables: ['userName', 'userDisplayName'],
	},
	STREAM_FIRST_USER_JOIN: {
		description: 'User joined stream for first time',
		availableVariables: ['userName'],
	},
	RAIDED: {
		description: '',
		availableVariables: ['userName', 'userDisplayName', 'raidViewers'],
	},
	TITLE_OR_CATEGORY_CHANGED: {
		description: '',
		availableVariables: [
			'oldStreamTitle',
			'newStreamTitle',
			'oldStreamCategory',
			'newStreamCategory',
		],
	},
	STREAM_ONLINE: {
		description: '',
		availableVariables: ['streamTitle', 'streamCategory'],
	},
	STREAM_OFFLINE: {
		description: '',
		availableVariables: [],
	},
	ON_CHAT_CLEAR: {
		description: '',
		availableVariables: [],
	},
	DONATE: {
		description: '',
		availableVariables: ['userName', 'donateAmount', 'donateCurrency', 'donateMessage'],
	},
	KEYWORD_MATCHED: {
		description: 'Keyword used',
		availableVariables: ['userName', 'userDisplayName', 'keywordName', 'keywordResponse'],
	},
	GREETING_SENDED: {
		description: 'Greeting sended in chat',
		availableVariables: ['userName', 'userDisplayName', 'greetingText'],
	},
	POLL_BEGIN: {
		description: '[POLL] Poll started',
		availableVariables: ['pollTitle', 'pollOptionsNames'],
	},
	POLL_PROGRESS: {
		description: '[POLL] Vote received',
		availableVariables: ['pollTitle', 'pollOptionsNames', 'pollTotalVotes'],
	},
	POLL_END: {
		description: '[POLL] Poll ended',
		availableVariables: [
			'pollTitle',
			'pollOptionsNames',
			'pollTotalVotes',
			'pollWinnerTitle',
			'pollWinnerBitsVotes',
			'pollWinnerChannelsPointsVotes',
			'pollWinnerTotalVotes',
		],
	},
	PREDICTION_BEGIN: {
		description: '[PREDICTION] Prediction started',
		availableVariables: ['predictionTitle', 'predictionOptionsNames'],
	},
	PREDICTION_PROGRESS: {
		description: '[PREDICTION] Prediction received',
		availableVariables: [
			'predictionTitle',
			'predictionOptionsNames',
			'predictionTotalChannelPoints',
		],
	},
	PREDICTION_LOCK: {
		description: '[PREDICTION] Prediction locked.',
		availableVariables: [
			'predictionTitle',
			'predictionOptionsNames',
			'predictionTotalChannelPoints',
		],
	},
	PREDICTION_END: {
		description: '[PREDICTION] Prediction ended',
		availableVariables: [
			'predictionTitle',
			'predictionOptionsNames',
			'predictionTotalChannelPoints',
			`predictionWinner.title`,
			`predictionWinner.totalUsers`,
			`predictionWinner.totalPoints`,
			`predictionWinner.topUsers`,
		],
	},
};
