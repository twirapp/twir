import { Entity, Schema } from 'redis-om';

export class ModerationSettings extends Entity {
  id: string;
  type: string;
  channelId: string;
  enabled: boolean;
  subscribers: boolean;
  vips: boolean;
  banTime: number;
  banMessage: string | null;
  warningMessage: string | null;
  checkClips: boolean | null;
  triggerLength: number | null;
  maxPercentage: number | null;
  blackListSentences: string[] | null;
}

export const moderationSettingsSchema = new Schema(ModerationSettings, {
  id: { type: 'string' },
  type: { type: 'string', indexed: true },
  channelId: { type: 'string', indexed: true },
  enabled: { type: 'boolean' },
  subscribers: { type: 'boolean' },
  vips: { type: 'boolean' },
  banTime: { type: 'number' },
  banMessage: { type: 'string' },
  warningMessage: { type: 'string' },
  checkClips: { type: 'boolean' },
  triggerLength: { type: 'number' },
  maxPercentage: { type: 'number' },
  blackListSentences: { type: 'string[]' },
}, {
  prefix: 'settings:moderation',
  indexedDefault: true,
});
