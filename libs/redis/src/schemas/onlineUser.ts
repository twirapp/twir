import { Entity, Schema } from 'redis-om';

class OnlineUser extends Entity {
  userId: string;
  userName: string;
  channelId: string;
}

export const onlineUserSchema = new Schema(OnlineUser, {
  userId: { type: 'string' },
  userName: { type: 'string' },
  channelId: { type: 'string' },
}, {
  prefix: 'onlineUsers',
  indexedDefault: true,
});
