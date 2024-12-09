import { newBus } from '@twir/bus-core'
import _ from 'lodash'
import { connect as natsConnect } from 'nats'
import { VM } from 'vm2'

const vm = new VM({
	sandbox: {
		fetch,
		URLSearchParams,
		_,
	},
	timeout: 1000,
	wasm: false,
	eval: false,
})

const nc = await natsConnect({
	servers: process.env.NODE_ENV === 'production' ? 'nats://nats:4222' : 'nats://127.0.0.1:4222',
})
const bus = newBus(nc)

bus.Eval.Evaluate.subscribeGroup('eval.evaluate', async (request) => {
	let resultOfExecution
	try {
		const toEval = `(async function () { ${request.expression} })()`.split(';\n').join(';')
		resultOfExecution = await vm.run(toEval)
	} catch (error) {
		console.error(error)
		resultOfExecution = error.message ?? 'unexpected error'
	}

	return {
		result: String(resultOfExecution).slice(0, 5000),
	}
})

console.info('Eval service started')
