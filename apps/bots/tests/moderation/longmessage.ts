import { expect, test, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = {
  enabled: true,
  vips: true,
  subscribers: true,
  triggerLength: 10,
  banMessage: 'ban',
  banTime: 1,
  warningMessage: 'warning',
  blackListSentences: ['test'],
};
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  longMessage: settings,
}) as any);

test('Should moderate "qqqqqqqqqqq" (11)', async () => {
  const state = createState();
  const result = await parser.parse('qqqqqqqqqqq', state);

  expect(result?.time).toBe(1);
});

test('Should not moderate "qqqqqqqqq" (9)', async () => {
  const state = createState();
  const result = await parser.parse('qqqqqqqqq', state);

  expect(result).toBe(undefined);
});
