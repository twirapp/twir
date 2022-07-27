import { Entity, Schema } from 'redis-om';

export class Stream extends Entity {
  id: string;
  user_id: string;
  user_login: string;
  user_name: string;
  game_id: string;
  game_name: string;
  type: string;
  title: string;
  viewer_count: number;
  started_at: string;
  language: string;
  thumbnail_url: string;
  tag_ids: string[] | null;
  is_mature: boolean;
  parsedMessages?: number;
}

export const streamSchema = new Schema(Stream, {
  id: { type: 'string', indexed: true },
  user_id: { type: 'string', indexed: true },
  user_login: { type: 'string' },
  user_name: { type: 'string' },
  game_id: { type: 'string' },
  game_name: { type: 'string' },
  type: { type: 'string' },
  title: { type: 'string' },
  viewer_count: { type: 'number' },
  started_at: { type: 'string' },
  language: { type: 'string' },
  thumbnail_url: { type: 'string' },
  tag_ids: { type: 'string[]' },
  is_mature: { type: 'boolean' },
  parsedMessages: { type: 'number' },
}, {
  prefix: 'streams',
  indexedDefault: true,
});
