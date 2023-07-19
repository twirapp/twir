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

export const EVENTS: Record<
	string,
	{
		name: string;
		icon?: FunctionalComponent;
	}
> = {
	FOLLOW: {
		name: 'Follow',
		icon: IconUserHeart,
	},
	SUBSCRIBE: {
		name: 'Subscribe',
		icon: IconUserDollar,
	},
	RESUBSCRIBE: {
		name: 'Resubscribe',
		icon: IconUserStar,
	},
	SUB_GIFT: {
		name: 'Subscribe Gift',
		icon: IconGift,
	},
	REDEMPTION_CREATED: {
		name: 'Reward Activated',
		icon: IconAward,
	},
	COMMAND_USED: {
		name: 'Reward Activated',
		icon: IconStar,
	},
	FIRST_USER_MESSAGE: {
		name: 'First User Message',
		icon: IconMessageExclamation,
	},
	STREAM_FIRST_USER_JOIN: {
		name: 'User Joined Stream For First Time',
		icon: IconUserPlus,
	},
	RAIDED: {
		name: 'Raided',
		icon: IconPick,
	},
	TITLE_OR_CATEGORY_CHANGED: {
		name: 'Title or Category Changed',
		icon: IconTransform,
	},
	STREAM_ONLINE: {
		name: 'Stream Online',
		icon: IconAccessPoint,
	},
	STREAM_OFFLINE: {
		name: 'Stream Offline',
		icon: IconAccessPointOff,
	},
	ON_CHAT_CLEAR: {
		name: 'On Chat Clear',
		icon: IconEraser,
	},
	DONATE: {
		name: 'Donate',
		icon: IconCashBanknote,
	},
	KEYWORD_MATCHED: {
		name: 'Keyword Matched',
		icon: IconBracketsContain,
	},
	GREETING_SENDED: {
		name: 'Greeting Sended',
		icon: IconHeartHandshake,
	},
	POLL_BEGIN: {
		name: 'Poll Begin',
		icon: IconDeviceDesktopAnalytics,
	},
	POLL_PROGRESS: {
		name: 'Poll Progress',
		icon: IconDeviceDesktopAnalytics,
	},
	POLL_END: {
		name: 'Poll End',
		icon: IconDeviceDesktopAnalytics,
	},
	PREDICTION_BEGIN: {
		name: 'Prediction Begin',
		icon: IconDice6,
	},
	PREDICTION_PROGRESS: {
		name: 'Prediction Progress',
		icon: IconDice6,
	},
	PREDICTION_LOCK: {
		name: 'Prediction Lock',
		icon: IconDice6,
	},
	PREDICTION_END: {
		name: 'Prediction End',
		icon: IconDice6,
	},
};
