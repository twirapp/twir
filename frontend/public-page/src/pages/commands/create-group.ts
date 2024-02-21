import type { Command } from '@twir/api/messages/commands_unprotected/commands_unprotected';

export type Group = {
	name: string
	commands: Command[]
}

export const createGroups = (commands: Command[]): Array<Command | Group> => {
	const groups: Map<string, Group> = new Map();
	const result: Array<Command | Group> = [];

	for (const command of commands) {
		if (!command.group && command.module === 'CUSTOM') {
			result.push(command);
			continue;
		}

		const module = command.group ?? command.module!;
		if (!groups.get(module)) {
			groups.set(module, { name: module, commands: [] });
		}

		groups.get(module)!.commands.push(command);
	}

	return [...result, ...groups.values()];
};

export const isCommand = (item: Command | Group): item is Command => {
	return !('commands' in item);
};
