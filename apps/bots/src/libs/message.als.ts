import { AsyncLocalStorage } from 'async_hooks';

import { Logger } from '@nestjs/common';

import { ConsoleLogger } from './logger.js';

export const messageAls = new AsyncLocalStorage<{
  messageId: string,
  logger: Logger,
}>();