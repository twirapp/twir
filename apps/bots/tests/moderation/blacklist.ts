import { describe, expect, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { prisma } from '../../src/libs/prisma.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = {
  enabled: true,
  vips: true,
  subscribers: true,
  banMessage: 'ban',
  banTime: 1,
  warningMessage: 'warning',
  blackListSentences: ['test'],
};
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  blacklists: settings,
}) as any);

describe('Regular user', (t) => {
  t('Should moderate "test"', async () => {
    const state = createState();
    const result = await parser.parse('test', state);

    expect(result?.time).toBe(1);
  });

  t('Should not moderate "qwe"', async () => {
    const state = createState();
    const result = await parser.parse('qwe', state);

    expect(result).toBe(undefined);
  });
});
