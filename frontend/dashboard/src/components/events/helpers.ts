import type { SelectGroupOption, SelectOption } from 'naive-ui';

import { EVENTS } from './events.js';
import { OPERATIONS, Operation } from './operations.js';

type SelectGeneric = {
	type?: 'group'
	name: string,
	childrens?: Record<string, SelectGeneric>
}

export const createSelectOptions = (values: Record<string, SelectGeneric>): (SelectOption | SelectGroupOption)[] => {
	return Object.entries(values)
		.map(([key, value]) => {
			const result: SelectOption | SelectGroupOption = {
				value: key,
				label: value.name,
			};

			if (value.type === 'group' && value.childrens) {
				result.key = value.name;
				result.type = 'group';
				result.children = Object.entries(value.childrens).map(([childKey, childValue]) => ({
					value: childKey,
					label: childValue.name,
				}));
			}

			return result;
		});
};

export const eventTypeSelectOptions = createSelectOptions(EVENTS);
export const operationTypeSelectOptions = createSelectOptions(OPERATIONS);

const createFlat = <T extends SelectGeneric>(values: Record<string, T>): Record<string, T> => {
	return Object.entries(values).reduce((acc, curr) => {
		if (curr[1].type === 'group' && curr[1].childrens) {
			Object.entries(curr[1].childrens)
				.forEach(([key, value]) => acc[key] = value as T);
			return acc;
		}

		acc[curr[0]] = curr[1];
		return acc;
	}, {} as Record<string, T>);
};

export const flatEvents = createFlat(EVENTS);
export const flatOperations = createFlat(OPERATIONS);

export const getOperation = (type: string): Operation | undefined => {
	return flatOperations[type];
};
