import { Command as CommandType, CommandModule, CommandPermission, CooldownType, Prisma } from '@tsuwari/prisma';
import { Entity, Schema } from 'redis-om';

type CMD = CommandType & { responses: string[] }

export class Command extends Entity implements CMD {
  id: string;
  name: string;
  cooldown: number | null;
  cooldownType: CooldownType;
  enabled: boolean;
  aliases: string[];
  description: string | null;
  visible: boolean;
  channelId: string;
  permission: CommandPermission;
  default: boolean;
  module: CommandModule;
  defaultName: string | null;
  responses: string[];
}

export const commandSchema = new Schema(Command, {
  id: { type: 'string', indexed: true },
  name: { type: 'string', indexed: true },
  cooldown: { type: 'number' },
  cooldownType: { type: 'string' },
  enabled: { type: 'boolean' },
  aliases: { type: 'string[]' },
  description: { type: 'string' },
  visible: { type: 'boolean' },
  channelId: { type: 'string' },
  permission: { type: 'string' },
  default: { type: 'boolean' },
  module: { type: 'string' },
  defaultName: { type: 'string' },
  responses: { type: 'string[]' },
}, {
  prefix: 'commands',
  indexedDefault: true,
});
