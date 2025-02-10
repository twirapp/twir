import { Queue } from './queue.js'

import type { NatsConnection } from 'nats'

export function newBus(nc: NatsConnection) {
	return {
		Eval: {
			Evaluate: new Queue<EvalRequest, EvalResponse>(nc, 'eval.evaluate'),
		},
	}
}

export interface EvalRequest {
	expression: string
}

export interface EvalResponse {
	result: string
}
