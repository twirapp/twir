import { Entity, Schema } from 'redis-om';

export class Keyword extends Entity {
  id: string;
  channelId: string;
  text: string;
  response: string;
  enabled: boolean;
  cooldown: number | null;
}

export const keywordsSchema = new Schema(Keyword, {
  id: { type: 'string' },
  channelId: { type: 'string' },
  text: { type: 'string' },
  response: { type: 'string' },
  enabled: { type: 'boolean' },
  cooldown: { type: 'number' },
}, {
  prefix: 'keywords',
  indexedDefault: true,
});
