import { randomInt } from 'crypto';

import { onlineUserSchema, RedisORMService } from '@tsuwari/redis';
import _ from 'lodash';

import { app } from '../../index.js';
import { Module } from '../index.js';

const redisOm = app.get(RedisORMService);
const repository = redisOm.fetchRepository(onlineUserSchema);

export const random: Module[] = [
  {
    key: 'random',
    description: 'Random number from N to N',
    example: 'random|1-40',
    handler: (key, state, params) => {
      if (!params) return '';
      const [from, to] = params.split('-').map(Number);

      if ([from, to].some((n) => typeof n !== 'number' || isNaN(n))) return '';
      return randomInt(from!, to!).toString();
    },
  },
  {
    key: 'random.online.user',
    description: 'Choose random online user',
    handler: async (_key, state) => {
      const users = await repository.search().where('channelId').equal(state.channelId).all();

      if (!users.length) return;
      const randomed = _.sample(users)!;

      return randomed.userName;
    },
  },
];
