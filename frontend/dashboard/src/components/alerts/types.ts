import { CreateRequest } from '@twir/grpc/generated/api/api/alerts';

export type EditableAlert = CreateRequest & { id?: string }
