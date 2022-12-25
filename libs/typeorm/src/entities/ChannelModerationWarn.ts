import { Column, Entity, PrimaryColumn } from 'typeorm';

import { type SettingsType } from './ChannelModerationSetting';

@Entity('channels_moderation_warnings', { schema: 'public' })
export class ChannelModerationWarn {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @Column('text', { name: 'userId' })
  userId: string;

  @Column('text')
  reason: SettingsType;
}
