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
  triggerLength: 5,
  banMessage: 'ban',
  banTime: 1,
  warningMessage: 'warning',
  blackListSentences: ['test'],
};
vi.spyOn(parser, 'getModerationSettings').mockImplementation(() => ({
  longMessage: settings,
}) as any);

test('Should moderate 4 emotes', async () => {
  const state = parseTwitchMessage(`@badge-info=;badges=vip/1,artist-badge/1;client-nonce=69f11804e5b281c04255679f9b590290;color=#456073;display-name=Bot_stop;emote-only=1;emotes=69:0-9,11-20,22-31,33-42;first-msg=0;flags=;id=8474058d-a3bd-4752-b205-ac9dc419620f;mod=0;room-id=128644134;subscriber=0;tmi-sent-ts=1654707695067;turbo=0;user-id=38019880;user-type= :bot_stop!bot_stop@bot_stop.tmi.twitch.tv PRIVMSG #sadisnamenya :BloodTrail BloodTrail BloodTrail BloodTrail`) as TwitchPrivateMessage;
  const result = await parser.parse('qqqqqqqqqqq', state);

  expect(result?.time).toBe(1);
});

test('Should not moderate 5 emotes', async () => {
  const state = parseTwitchMessage(`@badge-info=;badges=vip/1,artist-badge/1;client-nonce=ff2d9afbff39bbfddd258d2fdbaddb22;color=#456073;display-name=Bot_stop;emote-only=1;emotes=11934:0-9,11-20,22-31,33-42,44-53;first-msg=0;flags=;id=b98489d7-5e44-4851-a9d8-bd094212fdf5;mod=0;room-id=128644134;subscriber=0;tmi-sent-ts=1654707712236;turbo=0;user-id=38019880;user-type= :bot_stop!bot_stop@bot_stop.tmi.twitch.tv PRIVMSG #sadisnamenya :vlambeerYV vlambeerYV vlambeerYV vlambeerYV vlambeerYV`) as TwitchPrivateMessage;
  const result = await parser.parse('', state);

  expect(result).toBe(undefined);
});
