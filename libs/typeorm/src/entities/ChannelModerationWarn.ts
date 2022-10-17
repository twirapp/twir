/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { SettingsType } from './ChannelModerationSetting.js';
import { type User } from './User.js';

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
