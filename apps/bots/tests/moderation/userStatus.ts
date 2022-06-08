import { describe, expect, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = { enabled: true, vips: true, subscribers: true, banMessage: 'ban', banTime: 1, warningMessage: 'warning' };
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  links: settings,
}) as any);

describe('Subscriber', (t) => {
  t('Should moderate subscriber', async () => {
    const state = createState({ sub: true });
    const result = await parser.parse('vk.com', state);

    expect(result?.time).toBe(1);
  });

  t('Should not moderate subscriber', async () => {
    const state = createState({ sub: true });
    settings.subscribers = false;
    const result = await parser.parse('vk.com', state);

    expect(result).toBe(undefined);
    settings.subscribers = true;
  });
});

describe('Vip', (t) => {
  t('Should moderate vip', async () => {
    const state = createState({ vip: true });
    const result = await parser.parse('vk.com', state);

    expect(result?.time).toBe(1);
  });

  t('Should not moderate vip', async () => {
    const state = createState({ vip: true });
    settings.vips = false;
    const result = await parser.parse('vk.com', state);

    expect(result).toBe(undefined);
    settings.vips = true;
  });
});
