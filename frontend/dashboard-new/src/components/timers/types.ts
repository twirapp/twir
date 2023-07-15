import type { Timer } from '@twir/grpc/generated/api/api/timers';

export type EditableTimer = Omit<Timer, 'id' | 'channelId' | 'lastTriggerMessageNumber'> & { id?: string }
