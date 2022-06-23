import { Injectable, OnModuleInit } from '@nestjs/common';

import { ParserCache } from './cache.js';

export type State = {
  message?: string,
  channelId: string;
  sender?: {
    id?: string;
    name?: string;
  };
  cache: ParserCache;
};

type Handler = (key: string, state: State, params?: string | null, chatMessage?: string) => number | string | Promise<string | number | undefined> | undefined;

export interface Module {
  key: string;
  description?: string;
  example?: string;
  visible?: boolean;
  handler: Handler;
}

@Injectable()
export class VariablesParser implements OnModuleInit {
  vars: {
    [x: string]: Handler;
  } = {};
  readonly #regular = /\$\(([^)|]+)(?:\|([^)]+))?\)/g;

  onModuleInit() {
    setTimeout(async () => {
      const modules = await import('./modules/index.js');

      this.#registerModules(Object.values(modules).flat());
    }, 1000);
  }

  #registerModules(modules: Array<Module>, rewrite = false) {
    for (const module of modules) {
      if (this.vars[module.key] && !rewrite) {
        throw new Error(`Module ${module.key} already registered`);
      }

      this.vars[module.key] = module.handler;
    }
  }

  async parse(response: string, state: State, chatMessage?: string) {
    let result = '';
    const parts = response.split(this.#regular);

    for (let i = 0; i < parts.length + 2; i += 3) {
      result += parts[i];
      if (i + 1 < parts.length) {
        const key = parts[i + 1];
        const params = parts[i + 2];
        if (!key) continue;
        const newValue = this.vars[key];
        if (newValue === undefined) {
          result += `$(${key})`;
        } else if (typeof newValue === 'function') {
          const value = await newValue(key, state, params, chatMessage);
          result += typeof value !== 'undefined' ? value.toString() : `$(${key + params ? `|${params}` : ''})`;
        } else continue;
      }
    }

    return result;
  }
}
