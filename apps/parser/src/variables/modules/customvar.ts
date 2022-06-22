import { CustomVarType } from '@tsuwari/prisma';
import { VM } from 'vm2';

import { Module } from '../index.js';


const vm = new VM({
  sandbox: {
    fetch,
    URLSearchParams,
  },
  timeout: 5000,
});

export const customvar: Module = {
  key: 'customvar',
  description: 'Custom variable',
  visible: false,
  async handler(key, state, params) {
    if (!params) return;
    const vars = await state.cache.getCustomVars();
    const myVar = vars.find(v => v.name === params);

    if (!myVar) return;

    if (myVar.type === CustomVarType.SCRIPT && myVar.evalValue) {
      try {
        const opCounterFnc = 'let __opCount__ = 0; function __opCounter__() { if (__opCount__ > 10) { throw new Error("Running script seems to be in loop."); } else { __opCount__++; }};';
        const toEval = `(async function () { ${opCounterFnc} ${myVar.evalValue.split(';\n').map(line => '__opCounter__();' + line).join(';')} })`.replace(/\n/g, '');
        const result = await vm.run(toEval)();

        if (typeof result === 'string') return result;
        else return;
      } catch (e) {
        if (e instanceof Error) {
          return e.message;
        } else {
          return `Unexpected error happend on processing ${key}|${params} variable`;
        }
      }
    } else if (myVar.type === CustomVarType.TEXT && myVar.response) {
      return myVar.response;
    } else return;
  },
};