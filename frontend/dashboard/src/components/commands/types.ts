import type { Command, Command_Response } from '@twir/api/messages/commands/commands';

export type EditableCommand = Omit<
	Command,
	'responses' |
	'channelId' |
	'default' |
	'defaultName' |
	'id' |
	'group'
> & {
	responses: Array<Omit<Command_Response, 'id' | 'commandId' | 'order'>>,
	id?: string
};
