import { randomInt } from 'crypto';

import { Module } from '../index.js';

export const random: Module = {
  key: 'random',
  description: 'Random number from N to N',
  example: 'random|1-40',
  handler: (key, state, params) => {
    if (!params) return '';
    const [from, to] = params.split('-').map(Number);

    if ([from, to].some((n) => typeof n !== 'number' || isNaN(n))) return '';
    return randomInt(from!, to!).toString();
  },
};
