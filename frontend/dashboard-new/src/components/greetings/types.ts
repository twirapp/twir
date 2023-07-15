import type { CreateRequest } from '@twir/grpc/generated/api/api/greetings';

export type EditableGreeting = CreateRequest & {
	id?: string;
}
