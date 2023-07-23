import type { CreateData, CreateData_Response } from '@twir/grpc/generated/api/api/timers';

export type EditableTimerResponse = CreateData_Response

export type EditableTimer = CreateData & {
	id?: string,
}
