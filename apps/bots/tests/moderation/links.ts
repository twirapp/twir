/* import { expect, test, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = { enabled: true, vips: true, subscribers: true, banMessage: 'ban', banTime: 1, warningMessage: 'warning' };
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  links: settings,
}) as any);

test('Should moderate "vk . com"', async () => {
  const state = createState();
  const result = await parser.parse('vk . com', state);

  expect(result?.time).toBe(1);
});

test('Should moderate "vk. com"', async () => {
  const state = createState();
  const result = await parser.parse('vk. com', state);

  expect(result?.time).toBe(1);
});

test('Should moderate "vk .com"', async () => {
  const state = createState();
  const result = await parser.parse('vk .com', state);

  expect(result?.time).toBe(1);
});


test('Should moderate "qweqwe vk.com"', async () => {
  const state = createState();
  const result = await parser.parse('qweqwe vk.com', state);

  expect(result?.time).toBe(1);
});

test('Should not moderate "qweqwe"', async () => {
  const state = createState({ mod: true, sub: true, broadcaster: true });
  const result = await parser.parse('qweqwe', state);

  expect(result).toBe(undefined);
});

test('Should not moderate if permit', async () => {
  const state = createState();
  vi.spyOn(prisma.permit, 'findFirst').mockImplementation(() => Promise.resolve(true) as any);

  const result = await parser.parse('vk.com', state);
  expect(result).toBe(undefined);
});
 */
