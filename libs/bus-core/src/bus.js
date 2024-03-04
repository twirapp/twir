import { Queue } from './queue.js';

export const newBus = (nc) => {
	return {
		Eval: {
			Evaluate: new Queue(nc, 'eval.evaluate'),
		},
	};
};
