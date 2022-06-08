import { describe, expect, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { prisma } from '../../src/libs/prisma.js';
import { redis } from '../../src/libs/redis.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = { enabled: true, vips: true, subscribers: true, banMessage: 'ban', banTime: 1, warningMessage: 'warning' };
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  links: settings,
}) as any);

describe('Regular user', (t) => {
  t('Should moderate "vk . com"', async () => {
    const state = createState();
    const result = await parser.parse('vk . com', state);

    expect(result?.time).toBe(1);
  });

  t('Should moderate "vk. com"', async () => {
    const state = createState();
    const result = await parser.parse('vk. com', state);

    expect(result?.time).toBe(1);
  });

  t('Should moderate "vk .com"', async () => {
    const state = createState();
    const result = await parser.parse('vk .com', state);

    expect(result?.time).toBe(1);
  });


  t('Should moderate "qweqwe vk.com"', async () => {
    const state = createState();
    const result = await parser.parse('qweqwe vk.com', state);

    expect(result?.time).toBe(1);
  });

  t('Should be undefined on "qweqwe"', async () => {
    const state = createState({ mod: true, sub: true, broadcaster: true });
    const result = await parser.parse('qweqwe', state);

    expect(result).toBe(undefined);
  });
});

describe('Test Permit', (t) => {
  t('Should not moderate if permit', async () => {
    const state = createState();
    vi.spyOn(prisma.permit, 'findFirst').mockImplementation(() => Promise.resolve(true) as any);

    const result = await parser.parse('vk.com', state);
    expect(result).toBe(undefined);
  });
});
