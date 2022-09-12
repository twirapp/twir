import { config } from '@tsuwari/config';
import * as Eval from '@tsuwari/nats/eval';
import { connect } from 'nats';
import { VM } from 'vm2';

export const nats = await connect({
  servers: [config.NATS_URL],
});

const vm = new VM({
  sandbox: {
    fetch,
    URLSearchParams,
  },
  timeout: 1000,
});

(async () => {
  for await (const m of nats.subscribe('eval')) {
    const { script } = Eval.Evaluate.fromBinary(m.data);

    const opCounterFnc =
      'let __opCount__ = 0; function __opCounter__() { if (__opCount__ > 50000) { throw new Error("Running script seems to be in infinite loop."); } else { __opCount__++; }};';
    const toEval = `(async function () { ${opCounterFnc} ${script
      .split(';\n')
      .map((line) => '__opCounter__();' + line)
      .join(';')} })`.replace(/\n/g, '');
    const resultOfExecution = await vm.run(toEval)();

    const result = Eval.EvaluateResult.toBinary({
      result:
        typeof resultOfExecution === 'string' ? resultOfExecution : 'unexpected internal behavior.',
    });

    m.respond(result);
  }
})();
