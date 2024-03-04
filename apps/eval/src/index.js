import { newBus } from '@twir/bus-core';
import { config } from '@twir/config';
import _ from 'lodash';
import { connect as natsConnect } from 'nats';
import { VM } from 'vm2';

const vm = new VM({
	sandbox: {
		fetch,
		URLSearchParams,
		_: _,
	},
	timeout: 1000,
	wasm: false,
	eval: false,
});

const bus = newBus(await natsConnect({ servers: config.NATS_URL }));

bus.Eval.Evaluate.subscribeGroup('eval.evaluate', async (request) => {
	let resultOfExecution;
	try {
		const toEval = `(async function () { ${request.expression} })()`.split(';\n').join(';');
		resultOfExecution = await vm.run(toEval);
	} catch (error) {
		console.error(error);
		resultOfExecution = error.message ?? 'unexpected error';
	}

	return {
		result: String(resultOfExecution).slice(0, 5000),
	};
});


console.info('Eval service started');
