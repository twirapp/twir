import { type Command } from '@twir/api/messages/commands/commands';

export type Group = {
	name: string
	color?: string
	commands: Command[]
}

export const createGroups = (commands: Command[]): Array<Command | Group> => {
	const groups: Map<string, Group> = new Map();
	const result: Array<Command | Group> = [];

	for (const command of commands) {
		if (command.module !== 'CUSTOM') {
			const moduleName = command.module.toLowerCase();
			if (!groups.get(moduleName)) {
				groups.set(moduleName, { name: moduleName, commands: [] });
			}
			groups.get(moduleName)!.commands.push(command);
			continue;
		}

		if (!command.group || !command.groupId) {
			result.push(command);
			continue;
		}
		if (!groups.get(command.groupId)) {
			groups.set(command.groupId!, {
				name: command.group.name,
				color: command.group.color,
				commands: [],
			});
		}

		groups.get(command.groupId)!.commands.push(command);
	}

	return [...result, ...groups.values()];
};

export const isCommand = (item: Command | Group): item is Command => {
	return !('commands' in item);
};
