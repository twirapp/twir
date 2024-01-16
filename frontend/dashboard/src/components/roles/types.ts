import type { CreateRequest } from '@twir/api/messages/roles/roles';

export type EditableRole = CreateRequest & { id?: string };
