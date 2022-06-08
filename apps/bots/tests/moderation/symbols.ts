import { expect, test, vi } from 'vitest';

import './__mocks__/index.mock.js';

import { ModerationParser } from '../../src/libs/moderationParser.js';
import { createState } from './helpers.js';

const parser = new ModerationParser();

const settings = {
  enabled: true,
  vips: true,
  subscribers: true,
  maxPercentage: 50,
  banMessage: 'ban',
  banTime: 1,
  warningMessage: 'warning',
  blackListSentences: ['test'],
};
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  symbols: settings,
}) as any);

test('Should moderate "/////qqqqq"', async () => {
  const state = createState();
  const result = await parser.parse('/////qqqqq', state);

  expect(result?.time).toBe(1);
});

test('Should not moderate "qqqqqqqqqq"', async () => {
  const state = createState();
  const result = await parser.parse('qqqqqqqqqq', state);

  expect(result).toBe(undefined);
});
