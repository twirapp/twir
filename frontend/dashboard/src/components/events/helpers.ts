import { EventsOptions } from './events.js'
import { OPERATIONS } from './operations.js'

interface SelectGeneric {
	type?: 'group'
	name: string
	childrens?: Record<string, SelectGeneric>
}

function createFlat<T extends SelectGeneric>(values: Record<string, T>): Record<string, T> {
	return Object.entries(values).reduce((acc, curr) => {
		if (curr[1].type === 'group' && curr[1].childrens) {
			Object.entries(curr[1].childrens)
				.forEach(([key, value]) => acc[key] = value as T)
			return acc
		}

		acc[curr[0]] = curr[1]
		return acc
	}, {} as Record<string, T>)
}

export const flatEvents = createFlat(EventsOptions)
export const flatOperations = createFlat(OPERATIONS)

export const getEventName = (eventType: string) => flatEvents[eventType]?.name ?? eventType
