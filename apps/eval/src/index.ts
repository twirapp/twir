import { config } from '@tsuwari/config';
import * as Eval from '@tsuwari/nats/eval';
import _ from 'lodash';
import { connect } from 'nats';
import { VM } from 'vm2';

export const nats = await connect({
  servers: [config.NATS_URL],
});

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

const sub = nats.subscribe('eval');

(async () => {
  for await (const m of sub) {
    const { script } = Eval.Evaluate.fromBinary(m.data);

    let resultOfExecution: any;
    try {
      const toEval = `(async function () { ${script} })()`.replace(/\n/g, '');
      resultOfExecution = await vm.run(toEval);
    } catch (error) {
      resultOfExecution = (error as any).message ?? 'unexpected error';
    }

    const result = Eval.EvaluateResult.toBinary({
      result:
        typeof resultOfExecution === 'string' ? resultOfExecution : 'unexpected internal behavior.',
    });

    m.respond(result);
  }
})();
