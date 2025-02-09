import vm from 'node:vm'

import { newBus } from '@twir/bus-core'
import _ from 'lodash'
import { connect } from 'nats'

const nc = await connect({
	servers: Bun.env.NODE_ENV === 'production'
		? 'nats://nats:4222'
		: 'nats://127.0.0.1:4222',
})
const bus = newBus(nc)

bus.Eval.Evaluate.subscribeGroup('eval.evaluate', async (request) => {
	let resultOfExecution
	try {
		const toEval = `(async function () { ${request.expression} })()`
			.split(';\n')
			.join(';')

		const script = new vm.Script(toEval)
		const context = {
			fetch,
			URLSearchParams,
			_,
			result: undefined,
		}

		resultOfExecution = await script.runInContext(vm.createContext(context))
	} catch (error) {
		console.error(error)
		resultOfExecution = (error as any).message || 'unexpected error'
	}

	return {
		result: String(resultOfExecution).slice(0, 5000),
	}
})

console.info('Eval service started')
