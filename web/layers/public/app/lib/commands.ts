import type { PublicCommandsQuery } from '~/gql/graphql.js'

export type Command = PublicCommandsQuery['commandsPublic'][number]

export interface Group {
	name: string
	commands: Command[]
}

export function createGroups(commands: Command[]): Array<Command | Group> {
	const groups: Map<string, Group> = new Map()
	const result: Array<Command | Group> = []

	for (const command of commands) {
		if (!command.group && command.module === 'CUSTOM') {
			result.push(command)
			continue
		}

		const module = command.group?.name ?? command.module!
		if (!groups.get(module)) {
			groups.set(module, { name: module, commands: [] })
		}

		groups.get(module)!.commands.push(command)
	}

	return [...result, ...groups.values()]
}

export function isCommand(item: Command | Group): item is Command {
	return !('commands' in item)
}
