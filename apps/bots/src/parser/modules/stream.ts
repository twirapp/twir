import { humanizeStreamDuration } from '../../functions/humanizeStreamDuration.js';
import { Module } from '../index.js';

export const stream: Module[] = [
  {
    key: 'stream.title',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.title ?? '';
    },
  },
  {
    key: 'stream.uptime',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return !stream ? 'Stream offline.' : humanizeStreamDuration(new Date(stream.started_at).getTime());
    },
  },
  {
    key: 'stream.viewers',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.viewer_count ?? '';
    },
  },
  {
    key: 'stream.category',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.game_name ?? '';
    },
  },
  {
    key: 'stream.messages',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.parsedMessages ?? 0;
    },
  },
];
