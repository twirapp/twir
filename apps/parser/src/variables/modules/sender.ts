import { Module } from '../index.js';

export const sender: Module = {
  key: 'sender',
  description: 'Username of user, who sended message',
  handler: (_, state) => state.sender?.name ?? '',
};
