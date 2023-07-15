import type { CreateRequest } from '@twir/grpc/generated/api/api/roles';

export type EditableRole = CreateRequest & { id?: string };
