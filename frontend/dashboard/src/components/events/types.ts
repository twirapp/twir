import type { Event, Event_Operation } from '@twir/grpc/generated/api/api/events';

export type EditableEvent = Omit<Event, 'id' | 'channelId'> & { id?: string }
export type EventOperation = Event_Operation
