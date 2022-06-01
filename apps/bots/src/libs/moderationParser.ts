import { ModerationSettings, SettingsType } from '@tsuwari/prisma';
import { TwitchPrivateMessage } from '@twurple/chat/lib/commands/TwitchPrivateMessage';
import tlds from 'tlds' assert { type: 'json' };

import { prisma } from './prisma.js';
import { redis } from './redis.js';

const clipsRegexps = [/.*(clips.twitch.tv\/)(\w+)/g, /.*(www.twitch.tv\/\w+\/clip\/)(\w+)/g];
const urlRegexps = [
  new RegExp(`(www)? ??\\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\\. ?(${tlds.join('|')})(?=\\P{L}|$)`, 'igu'),
  new RegExp(`[a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})(?=\\P{L}|$)`, 'igu'),
];

// @TODO: update redis cache on changes from panel
export class ModerationParser {
  async #getModerationSettings(channelId: string) {
    const result = {} as Record<SettingsType, ModerationSettings>;
    const settingsKeys = Object.values(SettingsType);

    for (const key of settingsKeys) {
      const redisKey = `settings:moderation:${channelId}:${key}`;
      const cachedSettings = await redis.get(redisKey);

      if (cachedSettings) {
        result[key] = JSON.parse(cachedSettings) as ModerationSettings;
      } else {
        const entity = await prisma.moderationSettings.findFirst({ where: { channelId: channelId, type: key } });
        if (entity) {
          redis.set(redisKey, JSON.stringify(entity), 'EX', 5 * 60 * 60);
          result[key] = entity;
        }
      }
    }

    return result;
  }

  async parse(message: string, state: TwitchPrivateMessage) {
    if (!state?.channelId) return;
    const settings = await this.#getModerationSettings(state.channelId);

    const results = await Promise.all([
      this.#linksParser(message, settings.links, state),
    ]);

    return results.find(r => typeof r !== 'undefined');
  }

  async #linksParser(message: string, settings: ModerationSettings | null, state: TwitchPrivateMessage) {
    if (!settings || !settings.enabled) return;
    if (settings.vips && state.userInfo.isVip) return;
    if (settings.subscribers && state.userInfo.isSubscriber) return;

    const containLink = urlRegexps.some(r => r.test(message));
    if (!containLink) return;

    if (!settings.checkClips && clipsRegexps.some(r => r.test(message))) return;

    const permit = await prisma.permit.findFirst({ where: { channelId: state.channelId!, userId: state.userInfo.userId } });
    if (permit) {
      await prisma.permit.delete({ where: { id: permit.id } });
      return;
    }

    const redisKey = `moderation:warnings:links:${state.userInfo.userId}`;
    const isWarned = await redis.get(redisKey);

    if (isWarned !== null) {
      redis.del(redisKey);
      return {
        time: settings.banTime,
        message: settings.banMessage,
      };
    } else {
      redis.set(redisKey, '', 'EX', 60 * 60);
      return {
        message: settings.warningMessage,
        delete: true,
      };
    }
  }
}