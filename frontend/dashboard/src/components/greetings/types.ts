import type { CreateRequest } from '@twir/api/messages/greetings/greetings';

export type EditableGreeting = CreateRequest & {
	id?: string;
	userName?: string
}
