import type { Event } from '@twir/grpc/generated/api/api/events';

export type EditableEvent = Omit<Event, 'id'> & { id?: string }
