import type { CreateRequest } from '@twir/api/messages/variables/variables';

export type EditableVariable = CreateRequest & {
	id?: string,
};
