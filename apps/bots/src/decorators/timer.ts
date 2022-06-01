/* eslint-disable prefer-rest-params */
import pc from 'picocolors';

import { messageAls } from '../libs/message.als.js';

export function timer() {
  return (_target: any, key: string | symbol, descriptor: PropertyDescriptor) => {
    const method = descriptor.value;
    const msgAls = messageAls.getStore();
    if (method.constructor.name === 'AsyncFunction') {
      descriptor.value = async function () {
        const start = performance.now();

        const result = await method.apply(this, arguments);

        if (msgAls) {
          msgAls?.logger.log(`Timer | ${this.constructor.name}.${key.toString()}() ${pc.bold(
            (performance.now() - start).toFixed(2),
          )}ms`);
        }
        return result;
      };
    } else {
      descriptor.value = function () {
        const start = performance.now();
        const result = method.apply(this, arguments);
        msgAls?.logger.log(`Timer | ${this.constructor.name}.${key.toString()}() ${pc.bold(
          (performance.now() - start).toFixed(2),
        )}ms`);
        return result;
      };
    }
  };
}
