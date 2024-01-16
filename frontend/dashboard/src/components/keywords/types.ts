import type { CreateRequest } from '@twir/api/messages/keywords/keywords';


export type EditableKeyword = CreateRequest & { id?: string }
