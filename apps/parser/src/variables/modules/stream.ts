
import { intervalToDuration, format, formatDuration } from 'date-fns';

import { Module } from '../index.js';

function humanizeStreamDuration(start: number, locale = 'en') {
  const duration = intervalToDuration({ start: start, end: Date.now() });

  return formatDuration(duration);
}
export const stream: Module[] = [
  {
    key: 'stream.title',
    description: 'Stream title',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.title ?? '';
    },
  },
  {
    key: 'stream.uptime',
    description: 'stream.uptime',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return !stream ? 'Stream offline.' : humanizeStreamDuration(new Date(stream.started_at).getTime());
    },
  },
  {
    key: 'stream.viewers',
    description: 'Stream vieweirs',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.viewer_count ?? '';
    },
  },
  {
    key: 'stream.category',
    description: 'Stream category',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.game_name ?? '';
    },
  },
  {
    key: 'stream.messages',
    description: 'Messages sended by users in this stream',
    handler: async (_, state) => {
      const stream = await state.cache.getStream();
      return stream?.parsedMessages ?? 0;
    },
  },
];
