import type { CreateData, CreateData_Response } from '@twir/api/messages/timers/timers';

export type EditableTimerResponse = CreateData_Response

export type EditableTimer = CreateData & {
	id?: string,
}
