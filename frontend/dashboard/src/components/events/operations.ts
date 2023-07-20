export type Operation = {
	description: string;
	haveInput?: boolean;
	additionalValues?: Array<string>;
	producedVariables?: Array<string>;
	dependsOnEvents?: string[];
	color?: 'default' | 'success' | 'error' | 'warning' | 'info' | 'primary';
	type?: 'group'
	childrens?: Record<string, Operation>
}

export const OPERATIONS: Record<string, Operation> = {
	SEND_MESSAGE: {
		description: 'Send message in chat',
		haveInput: true,
		additionalValues: ['useAnnounce'],
		color: 'success',
	},

	BANS: {
		description: 'Bans, Timeouts',
		type: 'group',
		childrens: {
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
		},
	},

	VIPS: {
		description: 'Manage vips',
		type: 'group',
		childrens: {
			VIP: {
				description: 'Vip user',
				haveInput: true,
				color: 'info',
			},
			UNVIP: {
				description: 'Unvip user',
				haveInput: true,
				color: 'error',
			},
			UNVIP_RANDOM: {
				description: 'Unvip random user',
				producedVariables: ['unvipedUserName'],
				color: 'warning',
			},
			UNVIP_RANDOM_IF_NO_SLOTS: {
				description: 'Unvip random user if no slots',
				haveInput: true,
				producedVariables: ['unvipedUserName'],
				color: 'warning',
			},
		},
	},

	MODS: {
		description: 'Mamange moderators',
		type: 'group',
		childrens: {
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
		},
	},

	STREAM: {
		description: 'Manage stream',
		type: 'group',
		childrens: {
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
		description: 'Manage chat',
		type: 'group',
		childrens: {
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
		},
	},

	CREATE_GREETING: {
		description:
			'Create greeting for user. Available only for rewards event, and requires user input.',
		dependsOnEvents: ['REDEMPTION_CREATED'],
		color: 'info',
	},

	OBS: {
		description: 'OBS Manage',
		type: 'group',
		childrens: {
			OBS_SET_SCENE: {
				description: 'Change scene',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_TOGGLE_SOURCE: {
				description: `Toggle source visibility`,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_TOGGLE_AUDIO: {
				description: 'Toggle audio on/off',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_SET_VOLUME: {
				description: 'Set audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_DECREASE_VOLUME: {
				description: 'Decrease audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_AUDIO_INCREASE_VOLUME: {
				description: 'Increase audio volume',
				haveInput: true,
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_ENABLE_AUDIO: {
				description: 'Enable audio source',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_DISABLE_AUDIO: {
				description: 'Disable audio source',
				additionalValues: ['target'],
				color: 'default',
			},
			OBS_START_STREAM: {
				description: 'Start stream',
				color: 'default',
			},
			OBS_STOP_STREAM: {
				description: 'Stop stream',
				color: 'default',
			},
		},
	},

	VARIABLES: {
		description: 'Variables',
		type: 'group',
		childrens: {
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
		},
	},

	TTS: {
		description: 'Text to speech (TTS)',
		type: 'group',
		childrens: {
			TTS_SAY: {
				description: 'Say text',
				haveInput: true,
				color: 'info',
			},
			TTS_DISABLE: {
				description: 'Disable TTS',
				color: 'info',
			},
			TTS_ENABLE: {
				description: 'Enable TTS',
				color: 'info',
			},
			TTS_SKIP: {
				description: 'Skip current text',
				color: 'info',
			},
			TTS_SWITCH_AUTOREAD: {
				description: 'Switch autoread messages on/off',
				color: 'info',
			},
			TTS_DISABLE_AUTOREAD: {
				description: 'Disable autoread messages',
				color: 'info',
			},
			TTS_ENABLE_AUTOREAD: {
				description: 'Enable autoread messages',
				color: 'info',
			},
		},
	},

	COMMANDS: {
		description: 'Commands manage',
		type: 'group',
		childrens: {
			ALLOW_COMMAND_TO_USER: {
				description: 'Allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',
			},
			REMOVE_ALLOW_COMMAND_TO_USER: {
				description: 'Remove allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',
			},
			DENY_COMMAND_TO_USER: {
				description: 'Deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',
			},
			REMOVE_DENY_COMMAND_TO_USER: {
				description: 'Remove deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',
			},
		},
	},
};
