import type { CreateRequest } from '@twir/grpc/generated/api/api/keywords';


export type EditableKeyword = CreateRequest & { id?: string }
