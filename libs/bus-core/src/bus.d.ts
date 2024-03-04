import type { NatsConnection } from 'nats';

import type { Queue } from './queue.d.ts';

export const newBus: (nc: NatsConnection) => {
	Eval: {
		Evaluate: Queue<EvalRequest, EvalResponse>;
	};
};


export interface EvalRequest {
	expression: string
}

export interface EvalResponse {
	result: string
}
