/* eslint-disable prefer-rest-params */
import { Logger } from '@nestjs/common';
import pc from 'picocolors';

const logger = new Logger('Timer');



const log = (name: string, key: string | symbol, time: number) => logger.log(`Timer | ${name}.${key.toString()}() ${pc.bold(
  (time).toFixed(2),
)}ms`);

export function timer() {
  return (_target: any, key: string | symbol, descriptor: PropertyDescriptor) => {
    const method = descriptor.value;
    if (method.constructor.name === 'AsyncFunction') {
      descriptor.value = async function () {
        const start = performance.now();

        const result = await method.apply(this, arguments);

        log(this.constructor.name, key, performance.now() - start);
        return result;
      };
    } else {
      descriptor.value = function () {
        const start = performance.now();
        const result = method.apply(this, arguments);
        log(this.constructor.name, key, performance.now() - start);
        return result;
      };
    }
  };
}
