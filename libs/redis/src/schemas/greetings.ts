import { Entity, Schema } from 'redis-om';

export class Greetings extends Entity {
  id: string;
  channelId: string;
  userId: string;
  enabled: boolean;
  text: string;
  processed: boolean | null;
}

export const greetingsSchema = new Schema(Greetings, {
  id: { type: 'string' },
  channelId: { type: 'string' },
  userId: { type: 'string' },
  enabled: { type: 'boolean' },
  text: { type: 'string' },
  processed: { type: 'boolean' },
}, {
  prefix: 'greetings',
  indexedDefault: true,
});
