import type { Command } from '~/gql/graphql'

export interface Group {
	name: string
	color?: string
	commands: Command[]
}

export function createGroups(commands: Command[]): Array<Command | Group> {
	const groups: Map<string, Group> = new Map()
	const result: Array<Command | Group> = []

	for (const command of commands) {
		if (command.module !== 'CUSTOM') {
			const moduleName = command.module.toLowerCase()
			if (!groups.get(moduleName)) {
				groups.set(moduleName, { name: moduleName, commands: [] })
			}
			groups.get(moduleName)!.commands.push(command)
			continue
		}

		if (!command.group) {
			result.push(command)
			continue
		}
		if (!groups.get(command.group.id)) {
			groups.set(command.group.id, {
				name: command.group.name,
				color: command.group.color,
				commands: []
			})
		}

		groups.get(command.group.id)!.commands.push(command)
	}

	return [...result, ...groups.values()]
}

export function isCommand(item: Command | Group): item is Command {
	return !('commands' in item)
}
