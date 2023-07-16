import type { Command, Command_Response } from '@twir/grpc/generated/api/api/commands';

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

export type ListRowData = Command & {
	isGroup?: boolean,
	groupColor?: string,
	children?: ListRowData[],
};
