import { CustomVarType } from '@tsuwari/prisma';
import { Entity, Schema } from 'redis-om';

export class CustomVar extends Entity {
  id: string;
  name: string;
  description: string | null;
  type: CustomVarType;
  evalValue: string | null;
  response: string | null;
  channelId: string;
}

export const customVarSchema = new Schema(CustomVar, {
  id: { type: 'string' },
  name: { type: 'string' },
  description: { type: 'string' },
  type: { type: 'string' },
  evalValue: { type: 'string' },
  response: { type: 'string' },
  channelId: { type: 'string' },
}, {
  prefix: 'variables',
  indexedDefault: true,
});
