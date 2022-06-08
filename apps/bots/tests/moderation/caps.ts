import { parseTwitchMessage } from '@twurple/chat';
import type { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';
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
  caps: settings,
}) as any);

test('Should not moderate', async () => {
  const state = createState();
  const result = await parser.parse('QQQQqqqqqq', state);

  expect(result).toBe(undefined);
});

test('Should moderate', async () => {
  const state = createState();
  const result = await parser.parse('QQQQQqqqqq', state);

  expect(result?.time).toBe(1);
});

test('Should not moderate "BloodTrail BloodTrail BloodTrail чел, а это сообщение капс??? BloodTrail BloodTrail BloodTrail ну, НАЙС БАН!!! BloodTrail BloodTrail BloodTrail"', async () => {
  const state = parseTwitchMessage(`@badge-info=;badges=vip/1,artist-badge/1;client-nonce=47e2ab31537d1cb66358fd57392df65e;color=#456073;display-name=Bot_stop;emotes=69:0-9,11-20,22-31,62-71,73-82,84-93,111-120,122-131,133-142;first-msg=0;flags=;id=78e7cc5e-0fe1-4179-91cd-6bace43bf707;mod=0;room-id=128644134;subscriber=0;tmi-sent-ts=1654706610325;turbo=0;user-id=38019880;user-type= :bot_stop!bot_stop@bot_stop.tmi.twitch.tv PRIVMSG #sadisnamenya :BloodTrail BloodTrail BloodTrail чел, а это сообщение капс??? BloodTrail BloodTrail BloodTrail ну, НАЙС БАН!!! BloodTrail BloodTrail BloodTrail`) as TwitchPrivateMessage;

  const result = await parser.parse('BloodTrail BloodTrail BloodTrail чел, а это сообщение капс??? BloodTrail BloodTrail BloodTrail ну, НАЙС БАН!!! BloodTrail BloodTrail BloodTrail', state);

  expect(result).toBe(undefined);
});
