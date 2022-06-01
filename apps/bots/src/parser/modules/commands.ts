import {
  getChannelCommandsByNamesFromRedis,
  getChannelCommandsNamesFromRedis,
} from '../../functions/getChannelCommandListFromRedis.js';
import { Module } from '../index.js';

export const commands: Module[] = [
  {
    key: 'commands.list',
    handler: async (_, state) => {
      const names = await getChannelCommandsNamesFromRedis(state.channelId) ?? [];
      const commands = await getChannelCommandsByNamesFromRedis(state.channelId, names);
      return commands.map((c) => `!${c.name}`).join(', ') ?? '';
    },
  },
];
