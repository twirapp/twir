import type { CreateRequest } from '@twir/grpc/generated/api/api/variables';

export type EditableVariable = CreateRequest & {
	id?: string,
};
