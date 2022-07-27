import { HelpersService } from '../../helpers.service.js';
import { app } from '../../index.js';
import { Module } from '../index.js';

const helpers = app.get(HelpersService);

export const commands: Module[] = [
  {
    key: 'commands.list',
    description: 'List of commands',
    handler: async (_, state) => {
      const commands = await helpers.getChannelCommands(state.channelId);
      const filteredCommands = commands.filter(c => c.visible ?? true).map((c) => `!${c.name}`).join(', ') ?? '';

      return filteredCommands;
    },
  },
];
