import { Module } from '../index.js';

export const sender: Module = {
  key: 'sender',
  handler: (_, state) => state.sender.name ?? '',
};
