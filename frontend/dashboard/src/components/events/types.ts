import type { Event, Event_Operation } from '@twir/api/messages/events/events';

export type EditableEvent = Omit<Event, 'id' | 'channelId' | 'message'> & { id?: string }
export type EventOperation = Event_Operation
