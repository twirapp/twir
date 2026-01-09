import { EventOperationType } from '~/gql/graphql.ts'

export interface Operation {
	name: string
	haveInput?: boolean
	additionalValues?: Array<string>
	producedVariables?: Array<string>
	dependsOnEvents?: string[]
	color?: 'default' | 'success' | 'error' | 'warning' | 'info'
	type?: 'group'
	childrens?: Partial<Record<EventOperationType, Operation>>
	inputKeyTranslatePath?: string
}

export const EventOperations: Record<EventOperationType | string, Operation> = {
	[EventOperationType.SendMessage]: {
		name: 'Send message in chat',
		haveInput: true,
		additionalValues: ['useAnnounce'],
		color: 'success',
		inputKeyTranslatePath: 'events.operations.inputs.message',
	},
	[EventOperationType.MessageDelete]: {
		name: 'Delete message',
		color: 'error',
		dependsOnEvents: ['COMMAND_USED'],
	},

	[EventOperationType.TriggerAlert]: {
		name: 'Trigger alert',
		color: 'success',
	},

	BANS: {
		name: 'Bans, Timeouts',
		type: 'group',
		childrens: {
			[EventOperationType.Ban]: {
				name: 'Ban user',
				haveInput: true,
				additionalValues: ['timeoutMessage'],
				color: 'error',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.Unban]: {
				name: 'Unban user',
				haveInput: true,
				color: 'success',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.BanRandom]: {
				name: 'Ban random online user',
				producedVariables: ['bannedUserName'],
				additionalValues: ['timeoutMessage'],
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.Timeout]: {
				name: 'Timeout user',
				haveInput: true,
				additionalValues: ['timeoutTime', 'timeoutMessage'],
				color: 'error',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.TimeoutRandom]: {
				name: 'Timeout random online user',
				producedVariables: ['bannedUserName'],
				additionalValues: ['timeoutTime', 'timeoutMessage'],
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
		},
	},

	VIPS: {
		name: 'Manage vips',
		type: 'group',
		childrens: {
			[EventOperationType.Vip]: {
				name: 'Vip user',
				haveInput: true,
				color: 'info',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.Unvip]: {
				name: 'Unvip user',
				haveInput: true,
				color: 'error',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.UnvipRandom]: {
				name: 'Unvip random user',
				producedVariables: ['unvipedUserName'],
				color: 'warning',
			},
			[EventOperationType.UnvipRandomIfNoSlots]: {
				name: 'Unvip random user if no slots',
				haveInput: true,
				producedVariables: ['unvipedUserName'],
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.vipSlots',
			},
		},
	},

	MODS: {
		name: 'Manage moderators',
		type: 'group',
		childrens: {
			[EventOperationType.Mod]: {
				name: 'Give user moderation',
				haveInput: true,
				color: 'info',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.Unmod]: {
				name: 'Remove moderation from user',
				haveInput: true,
				color: 'error',
				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.UnmodRandom]: {
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
			[EventOperationType.ChangeTitle]: {
				name: 'Change title of stream',
				haveInput: true,
				color: 'warning',
			},
			[EventOperationType.ChangeCategory]: {
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
			[EventOperationType.EnableSubmode]: {
				name: 'Enable submode',
				color: 'success',
			},
			[EventOperationType.DisableSubmode]: {
				name: 'Disable submode',
				color: 'error',
			},
			[EventOperationType.EnableEmoteOnly]: {
				name: 'Enable emoteonly',
				color: 'success',
			},
			[EventOperationType.DisableEmoteOnly]: {
				name: 'Disable emoteonly',
				color: 'error',
			},
		},
	},

	[EventOperationType.CreateGreeting]: {
		name:
			'Create greeting for user. Available only for rewards event, and requires user input.',
		dependsOnEvents: ['REDEMPTION_CREATED'],
		color: 'info',
	},

	OBS: {
		name: 'OBS Manage',
		type: 'group',
		childrens: {
			[EventOperationType.ObsChangeScene]: {
				name: 'Change scene',
				color: 'default',
			},
			[EventOperationType.ObsToggleSource]: {
				name: `Toggle source visibility`,
				color: 'default',
			},
			[EventOperationType.ObsToggleAudio]: {
				name: 'Toggle audio on/off',
				color: 'default',
			},
			[EventOperationType.ObsSetAudioVolume]: {
				name: 'Set audio volume',
				haveInput: true,
				color: 'default',
			},
			[EventOperationType.ObsDecreaseAudioVolume]: {
				name: 'Decrease audio volume',
				haveInput: true,
				color: 'default',
			},
			[EventOperationType.ObsIncreaseAudioVolume]: {
				name: 'Increase audio volume',
				haveInput: true,
				color: 'default',
			},
			[EventOperationType.ObsEnableAudio]: {
				name: 'Enable audio source',
				color: 'default',
			},
			[EventOperationType.ObsDisableAudio]: {
				name: 'Disable audio source',
				color: 'default',
			},
			[EventOperationType.ObsStartStream]: {
				name: 'Start stream',
				color: 'default',
			},
			[EventOperationType.ObsStopStream]: {
				name: 'Stop stream',
				color: 'default',
			},
		},
	},

	VARIABLES: {
		name: 'Variables',
		type: 'group',
		childrens: {
			[EventOperationType.ChangeVariable]: {
				name: 'Change variable',
				haveInput: true,
				additionalValues: ['target'],
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.variableValue',
			},
			[EventOperationType.DecrementVariable]: {
				name: 'Decrement number variable',
				haveInput: true,
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.variableValue',
			},
			[EventOperationType.IncrementVariable]: {
				name: 'Increment number variable',
				haveInput: true,
				additionalValues: ['target'],
				color: 'warning',
				inputKeyTranslatePath: 'events.operations.inputs.variableValue',
			},
		},
	},

	SEVENTV: {
		name: '7tv',
		type: 'group',
		childrens: {
			[EventOperationType.SeventvAddEmote]: {
				name: 'Add emote',
				haveInput: true,
				color: 'success',
			},
			[EventOperationType.SeventvRemoveEmote]: {
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
			[EventOperationType.TtsSay]: {
				name: 'Say text',
				haveInput: true,
				color: 'info',
			},
			[EventOperationType.TtsDisable]: {
				name: 'Disable TTS',
				color: 'info',
			},
			[EventOperationType.TtsEnable]: {
				name: 'Enable TTS',
				color: 'info',
			},
			[EventOperationType.TtsSkip]: {
				name: 'Skip current text',
				color: 'info',
			},
			[EventOperationType.TtsSwitchAutoread]: {
				name: 'Switch autoread messages on/off',
				color: 'info',
			},
			[EventOperationType.TtsDisableAutoread]: {
				name: 'Disable autoread messages',
				color: 'info',
			},
			[EventOperationType.TtsEnableAutoread]: {
				name: 'Enable autoread messages',
				color: 'info',
			},
		},
	},

	RAIDS: {
		name: 'Raids',
		type: 'group',
		childrens: {
			[EventOperationType.RaidChannel]: {
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
			[EventOperationType.AllowCommandToUser]: {
				name: 'Allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',

				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.RemoveAllowCommandToUser]: {
				name: 'Remove allow command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',

				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.DenyCommandToUser]: {
				name: 'Deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'error',

				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
			[EventOperationType.RemoveDenyCommandToUser]: {
				name: 'Remove deny command to user',
				haveInput: true,
				additionalValues: ['target'],
				color: 'success',

				inputKeyTranslatePath: 'events.operations.inputs.username',
			},
		},
	},

	[EventOperationType.ShoutoutChannel]: {
		name: 'Shoutout channel',
		haveInput: true,
		color: 'info',
	},
}
