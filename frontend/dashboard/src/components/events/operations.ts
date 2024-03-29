export type Operation = {
	name: string;
	haveInput?: boolean;
	additionalValues?: Array<string>;
	producedVariables?: Array<string>;
	dependsOnEvents?: string[];
	color?: 'default' | 'success' | 'error' | 'warning' | 'info';
	type?: 'group'
	childrens?: Record<string, Operation>
}

export const OPERATIONS: Record<string, Operation> = {
	SEND_MESSAGE: {
		name: 'Send message in chat',
		haveInput: true,
		additionalValues: ['useAnnounce'],
		color: 'success',
	},

	TRIGGER_ALERT: {
		name: 'Trigger alert',
		additionalValues: ['target'],
		color: 'success',
	},

	BANS: {
		name: 'Bans, Timeouts',
		type: 'group',
		childrens: {
			BAN: {
				name: 'Ban user',
				haveInput: true,
				additionalValues: ['timeoutMessage'],
				color: 'error',
			},
			UNBAN: {
				name: 'Unban user',
				haveInput: true,
				color: 'success',
			},
			BAN_RANDOM: {
				name: 'Ban random online user',
				producedVariables: ['bannedUserName'],
				additionalValues: ['timeoutMessage'],
				color: 'warning',
			},
			TIMEOUT: {
				name: 'Timeout user',
				haveInput: true,
				additionalValues: ['timeoutTime', 'timeoutMessage'],
				color: 'error',
			},
			TIMEOUT_RANDOM: {
				name: 'Timeout random online user',
				producedVariables: ['bannedUserName'],
				additionalValues: ['timeoutTime', 'timeoutMessage'],
				color: 'warning',
			},
		},
	},

	VIPS: {
		name: 'Manage vips',
		type: 'group',
		childrens: {
			VIP: {
				name: 'Vip user',
				haveInput: true,
				color: 'info',
			},
			UNVIP: {
				name: 'Unvip user',
				haveInput: true,
				color: 'error',
			},
			UNVIP_RANDOM: {
				name: 'Unvip random user',
				producedVariables: ['unvipedUserName'],
				color: 'warning',
			},
			UNVIP_RANDOM_IF_NO_SLOTS: {
				name: 'Unvip random user if no slots',
				haveInput: true,
				producedVariables: ['unvipedUserName'],
				color: 'warning',
			},
		},
	},

	MODS: {
		name: 'Manage moderators',
		type: 'group',
		childrens: {
			MOD: {
				name: 'Give user moderation',
				haveInput: true,
				color: 'info',
			},
			UNMOD: {
				name: 'Remove moderation from user',
				haveInput: true,
				color: 'error',
			},
			UNMOD_RANDOM: {
				name: 'Remove moderation from random user',
				producedVariables: ['unmodedUserName'],
				color: 'warning',
			},
		},
	},

	STREAM: {
		name: 'Manage stream',
		type: 'group',
		childrens: {
			CHANGE_TITLE: {
				name: 'Change title of stream',
				haveInput: true,
				color: 'warning',
			},
			CHANGE_CATEGORY: {
				name: 'Change category of stream',
				haveInput: true,
				color: 'warning',
			},
		},
	},
	// FULFILL_REDEMPTION: {
	// 	description: 'Verify fulfillment of the reward',
	// 	color: 'info',
	// },
	// CANCEL_REDEMPTION: {
	// 	description: 'Cancel reward and back points to user',
	// 	color: 'error',
	// },


	CHAT: {
		name: 'Manage chat',
		type: 'group',
		childrens: {
			ENABLE_SUBMODE: {
				name: 'Enable submode',
				color: 'success',
			},
			DISABLE_SUBMODE: {
				name: 'Disable submode',
				color: 'error',
			},
			ENABLE_EMOTEONLY: {
				name: 'Enable emoteonly',
				color: 'success',
			},
			DISABLE_EMOTEONLY: {
				name: 'Disable emoteonly',
				color: 'error',
			},
		},
	},

	CREATE_GREETING: {
		name:
			'Create greeting for user. Available only for rewards event, and requires user input.',
		dependsOnEvents: ['REDEMPTION_CREATED'],
		color: 'info',
	},

	OBS: {
		name: 'OBS Manage',
		type: 'group',
		childrens: {
			OBS_SET_SCENE: {
				name: 'Change scene',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_TOGGLE_SOURCE: {
				name: `Toggle source visibility`,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_TOGGLE_AUDIO: {
				name: 'Toggle audio on/off',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_SET_VOLUME: {
				name: 'Set audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_DECREASE_VOLUME: {
				name: 'Decrease audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_INCREASE_VOLUME: {
				name: 'Increase audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_ENABLE_AUDIO: {
				name: 'Enable audio source',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_DISABLE_AUDIO: {
				name: 'Disable audio source',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_START_STREAM: {
				name: 'Start stream',
				color: 'default',
			},
			OBS_STOP_STREAM: {
				name: 'Stop stream',
				color: 'default',
			},
		},
	},

	VARIABLES: {
		name: 'Variables',
		type: 'group',
		childrens: {
			CHANGE_VARIABLE: {
				name: 'Change variable',
				haveInput: true,
				additionalValues: ['target'],
				color: 'warning',
			},
			DECREMENT_VARIABLE: {
				name: 'Decrement number variable',
				haveInput: true,
				color: 'warning',
			},
			INCREMENT_VARIABLE: {
				name: 'Increment number variable',
				haveInput: true,
				additionalValues: ['target'],
				color: 'warning',
			},
		},
	},

	SEVENTV: {
		name: '7tv',
		type: 'group',
		childrens: {
			SEVENTV_ADD_EMOTE: {
				name: 'Add emote',
				haveInput: true,
				color: 'success',
			},
			SEVENTV_REMOVE_EMOTE: {
				name: 'Remove emote',
				haveInput: true,
				color: 'error',
			},
		},
	},

	TTS: {
		name: 'Text to speech (TTS)',
		type: 'group',
		childrens: {
			TTS_SAY: {
				name: 'Say text',
				haveInput: true,
				color: 'info',
			},
			TTS_DISABLE: {
				name: 'Disable TTS',
				color: 'info',
			},
			TTS_ENABLE: {
				name: 'Enable TTS',
				color: 'info',
			},
			TTS_SKIP: {
				name: 'Skip current text',
				color: 'info',
			},
			TTS_SWITCH_AUTOREAD: {
				name: 'Switch autoread messages on/off',
				color: 'info',
			},
			TTS_DISABLE_AUTOREAD: {
				name: 'Disable autoread messages',
				color: 'info',
			},
			TTS_ENABLE_AUTOREAD: {
				name: 'Enable autoread messages',
				color: 'info',
			},
		},
	},

	RAIDS: {
		name: 'Raids',
		type: 'group',
		childrens: {
			RAID_CHANNEL: {
				name: 'Raid channel',
				haveInput: true,
				color: 'warning',
			},
		},
	},

	COMMANDS: {
		name: 'Commands manage',
		type: 'group',
		childrens: {
			ALLOW_COMMAND_TO_USER: {
				name: 'Allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',
			},
			REMOVE_ALLOW_COMMAND_TO_USER: {
				name: 'Remove allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',
			},
			DENY_COMMAND_TO_USER: {
				name: 'Deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',
			},
			REMOVE_DENY_COMMAND_TO_USER: {
				name: 'Remove deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',
			},
		},
	},

	SHOUTOUT_CHANNEL: {
		name: 'Shoutout channel',
		haveInput: true,
		color: 'info',
	},
};
