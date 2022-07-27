import { Entity, Schema } from 'redis-om';

export class UsersStats extends Entity {
  id: string;
  userId: string;
  channelId: string;
  messages: number;
  watched: string;
}

export const usersStatsSchema = new Schema(UsersStats, {
  id: { type: 'string' },
  userId: { type: 'string' },
  channelId: { type: 'string' },
  messages: { type: 'number' },
  watched: { type: 'string' },
}, {
  prefix: 'usersStats',
  indexedDefault: true,
});
