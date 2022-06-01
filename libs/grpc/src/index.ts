export * from './bots.js';
export * from './streamstatus.js';
export * from './watched.js';

import { resolve } from 'path';
import { fileURLToPath } from 'url';

export const resolveProtoPath = (proto: string) => {
  return resolve(fileURLToPath(import.meta.url), '..', '..', `${proto}.proto`);
};