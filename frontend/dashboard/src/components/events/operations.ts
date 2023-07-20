export const operationMapping: Record<
	string,
	{
		description: string;
		haveInput?: boolean;
		additionalValues?: Array<string>;
		producedVariables?: Array<string>;
		dependsOnEvents?: Array<EventType>;
		color?: string;
	}
> = {
	SEND_MESSAGE: {
		description: 'Send message in chat',
		haveInput: true,
		additionalValues: ['useAnnounce'],
		color: 'success',
	},
	BAN: {
		description: 'Ban user',
		haveInput: true,
		additionalValues: ['timeoutMessage'],
		color: 'error',
	},
	UNBAN: {
		description: 'Unban user',
		haveInput: true,
		color: 'success',
	},
	BAN_RANDOM: {
		description: 'Ban random online user',
		producedVariables: ['bannedUserName'],
		additionalValues: ['timeoutMessage'],
		color: 'warning',
	},
	VIP: {
		description: '[VIPS] Vip user',
		haveInput: true,
		color: 'info',
	},
	UNVIP: {
		description: '[VIPS] Unvip user',
		haveInput: true,
		color: 'error',
	},
	UNVIP_RANDOM: {
		description: '[VIPS] Unvip random user',
		producedVariables: ['unvipedUserName'],
		color: 'warning',
	},
	UNVIP_RANDOM_IF_NO_SLOTS: {
		description: '[VIPS] Unvip random user if no slots',
		haveInput: true,
		producedVariables: ['unvipedUserName'],
		color: 'warning',
	},
	MOD: {
		description: 'Give user moderation',
		haveInput: true,
		color: 'info',
	},
	UNMOD: {
		description: 'Remove moderation from user',
		haveInput: true,
		color: 'error',
	},
	UNMOD_RANDOM: {
		description: 'Remove moderation from random user',
		producedVariables: ['unmodedUserName'],
		color: 'warning',
	},
	CHANGE_TITLE: {
		description: 'Change title of stream',
		haveInput: true,
		color: 'warning',
	},
	CHANGE_CATEGORY: {
		description: 'Change category of stream',
		haveInput: true,
		color: 'warning',
	},
	FULFILL_REDEMPTION: {
		description: 'Verify fulfillment of the reward',
		color: 'info',
	},
	CANCEL_REDEMPTION: {
		description: 'Cancel reward and back points to user',
		color: 'error',
	},
	ENABLE_SUBMODE: {
		description: 'Enable submode',
		color: 'success',
	},
	DISABLE_SUBMODE: {
		description: 'Disable submode',
		color: 'error',
	},
	ENABLE_EMOTEONLY: {
		description: 'Enable emoteonly',
		color: 'success',
	},
	DISABLE_EMOTEONLY: {
		description: 'Disable emoteonly',
		color: 'error',
	},
	CREATE_GREETING: {
		description:
			'Create greeting for user. Available only for rewards event, and requires user input.',
		dependsOnEvents: [EventType.REDEMPTION_CREATED],
		color: 'info',
	},
	TIMEOUT: {
		description: 'Timeout user',
		haveInput: true,
		additionalValues: ['timeoutTime', 'timeoutMessage'],
		color: 'error',
	},
	TIMEOUT_RANDOM: {
		description: 'Timeout random online user',
		producedVariables: ['bannedUserName'],
		additionalValues: ['timeoutTime', 'timeoutMessage'],
		color: 'warning',
	},
	OBS_SET_SCENE: {
		description: '[OBS] Change scene',
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_TOGGLE_SOURCE: {
		description: `[OBS] Toggle source visibility`,
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_TOGGLE_AUDIO: {
		description: '[OBS] Toggle audio on/off',
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_AUDIO_SET_VOLUME: {
		description: '[OBS] Set audio volume',
		haveInput: true,
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_AUDIO_DECREASE_VOLUME: {
		description: '[OBS] Decrease audio volume',
		haveInput: true,
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_AUDIO_INCREASE_VOLUME: {
		description: '[OBS] Increase audio volume',
		haveInput: true,
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_ENABLE_AUDIO: {
		description: '[OBS] Enable audio source',
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_DISABLE_AUDIO: {
		description: '[OBS] Disable audio source',
		additionalValues: ['target'],
		color: 'dark',
	},
	OBS_START_STREAM: {
		description: '[OBS] Start stream',
		color: 'dark',
	},
	OBS_STOP_STREAM: {
		description: '[OBS] Stop stream',
		color: 'dark',
	},
	CHANGE_VARIABLE: {
		description: 'Change variable',
		haveInput: true,
		additionalValues: ['target'],
		color: 'warning',
	},
	DECREMENT_VARIABLE: {
		description: 'Decrement number variable',
		haveInput: true,
		color: 'warning',
	},
	INCREMENT_VARIABLE: {
		description: 'Increment number variable',
		haveInput: true,
		additionalValues: ['target'],
		color: 'warning',
	},
	TTS_SAY: {
		description: '[TTS] Say text',
		haveInput: true,
		color: 'light',
	},
	TTS_DISABLE: {
		description: '[TTS] Disable TTS',
		color: 'light',
	},
	TTS_ENABLE: {
		description: '[TTS] Enable TTS',
		color: 'light',
	},
	TTS_SKIP: {
		description: '[TTS] Skip current text',
		color: 'light',
	},
	TTS_SWITCH_AUTOREAD: {
		description: '[TTS] Switch autoread messages on/off',
		color: 'light',
	},
	TTS_DISABLE_AUTOREAD: {
		description: '[TTS] Disable autoread messages',
		color: 'light',
	},
	TTS_ENABLE_AUTOREAD: {
		description: '[TTS] Enable autoread messages',
		color: 'light',
	},
	ALLOW_COMMAND_TO_USER: {
		description: '[COMMANDS] Allow command to user',
		haveInput: true,
		additionalValues: ['target'],
		color: 'success',
	},
	REMOVE_ALLOW_COMMAND_TO_USER: {
		description: '[COMMANDS] Remove allow command to user',
		haveInput: true,
		additionalValues: ['target'],
		color: 'error',
	},
	DENY_COMMAND_TO_USER: {
		description: '[COMMANDS] Deny command to user',
		haveInput: true,
		additionalValues: ['target'],
		color: 'error',
	},
	REMOVE_DENY_COMMAND_TO_USER: {
		description: '[COMMANDS] Remove deny command to user',
		haveInput: true,
		additionalValues: ['target'],
		color: 'success',
	},
};
