import { Entity, Schema } from 'redis-om';

class Message extends Entity {
  messageId: string;
  userId: string;
  channelId: string;
  userName: string;
  message: string;
  canBeDeleted: boolean;
}

export const messageSchema = new Schema(Message, {
  messageId: { type: 'string' },
  userId: { type: 'string' },
  channelId: { type: 'string' },
  userName: { type: 'string' },
  message: { type: 'text' },
  canBeDeleted: { type: 'boolean' },
}, {
  prefix: 'messages',
  indexedDefault: true,
});
