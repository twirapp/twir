/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, OneToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Bot } from './Bot.js';
import { type ChannelCommand } from './ChannelCommand.js';
import { type ChannelCustomvar } from './ChannelCustomvar.js';
import { type ChannelDotaAccount } from './ChannelDotaAccount.js';
import { type ChannelGreeting } from './ChannelGreeting.js';
import { type ChannelIntegration } from './ChannelIntegration.js';
import { type ChannelKeyword } from './ChannelKeyword.js';
import { type ChannelModerationSetting } from './ChannelModerationSetting.js';
import { type ChannelPermit } from './ChannelPermit.js';
import { type ChannelTimer } from './ChannelTimer.js';
import { type DashboardAccess } from './DashboardAccess.js';
import { type User } from './User.js';
import { type UserStats } from './UserStats.js';

@Entity('channels', { schema: 'public' })
export class Channel {
  @PrimaryColumn('text', { primary: true, name: 'id', unique: true })
  id: string;

  @Column('boolean', { name: 'isEnabled', default: true })
  isEnabled: boolean;

  @Column('boolean', { name: 'isTwitchBanned', default: false })
  isTwitchBanned: boolean;

  @Column('boolean', { name: 'isBanned', default: false })
  isBanned: boolean;

  @ManyToOne('Bot', 'channels', { onDelete: 'RESTRICT', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'botId', referencedColumnName: 'id' }])
  bot?: Relation<Bot>;

  @Column()
  botId: string;

  @OneToOne('User', 'channel', { onDelete: 'RESTRICT', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'id', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @OneToMany('ChannelCommand', 'channel')
  commands?: Relation<ChannelCommand[]>;

  @OneToMany('ChannelCustomvar', 'channel')
  customVar?: Relation<ChannelCustomvar[]>;

  @OneToMany('DashboardAccess', 'channel')
  dashboardAccess?: Relation<DashboardAccess[]>;

  @OneToMany('ChannelDotaAccount', 'channel')
  dotaAccounts?: Relation<ChannelDotaAccount[]>;

  @OneToMany('ChannelGreeting', 'channel')
  greetings?: Relation<ChannelGreeting[]>;

  @OneToMany('ChannelIntegration', 'channel')
  channelsIntegrations?: Relation<ChannelIntegration[]>;

  @OneToMany('ChannelKeyword', 'channel')
  keywords?: Relation<ChannelKeyword[]>;

  @OneToMany('ChannelModerationSetting', 'channel')
  moderationSettings?: Relation<ChannelModerationSetting[]>;

  @OneToMany('ChannelPermit', 'channel')
  permits?: Relation<ChannelPermit[]>;

  @OneToMany('ChannelTimer', 'channel')
  timers?: Relation<ChannelTimer[]>;

  @OneToMany('UserStats', 'channel')
  usersStats?: Relation<UserStats[]>;
}
