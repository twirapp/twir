import { CreateRequest } from '@twir/api/messages/alerts/alerts';

export type EditableAlert = CreateRequest & { id?: string }
