/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, OneToOne } from 'typeorm';

import { Bot } from './Bot.js';
import { ChannelCommand } from './ChannelCommand.js';
import { ChannelCustomvar } from './ChannelCustomvar.js';
import { DotaAccount } from './ChannelDotaAccount.js';
import { ChannelGreeting } from './ChannelGreeting.js';
import { ChannelIntegration } from './ChannelIntegration.js';
import { ChannelKeyword } from './ChannelKeyword.js';
import { ChannelModerationSetting } from './ChannelModerationSetting.js';
import { ChannelPermit } from './ChannelPermit.js';
import { ChannelTimer } from './ChannelTimer.js';
import { DashboardAccess } from './DashboardAccess.js';
import { User } from './User.js';
import { UserStats } from './UserStats.js';

@Index('channels_pkey', ['id'], { unique: true })
@Entity('channels', { schema: 'public' })
export class Channel {
  @Column('text', { primary: true, name: 'id' })
  id: string;

  @Column('boolean', { name: 'isEnabled', default: true })
  isEnabled: boolean;

  @Column('boolean', { name: 'isTwitchBanned', default: false })
  isTwitchBanned: boolean;

  @Column('boolean', { name: 'isBanned', default: false })
  isBanned: boolean;

  @ManyToOne(() => Bot, (bot) => bot.channels, { onDelete: 'RESTRICT', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'botId', referencedColumnName: 'id' }])
  bot: Bot;

  @OneToOne(() => User, (user) => user.channel, { onDelete: 'RESTRICT', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'id', referencedColumnName: 'id' }])
  user: User;

  @OneToMany(() => ChannelCommand, (command) => command.channel)
  commands: ChannelCommand[];

  @OneToMany(() => ChannelCustomvar, (customVar) => customVar.channel)
  customVar: ChannelCustomvar[];

  @OneToMany(() => DashboardAccess, (dashboardAccess) => dashboardAccess.channel)
  dashboardAccess: DashboardAccess[];

  @OneToMany(() => DotaAccount, (dotaAccount) => dotaAccount.channel)
  dotaAccounts: DotaAccount[];

  @OneToMany(() => ChannelGreeting, (greeting) => greeting.channel)
  greetings: ChannelGreeting[];

  @OneToMany(() => ChannelIntegration, (channelsIntegrations) => channelsIntegrations.channel)
  channelsIntegrations: ChannelIntegration[];

  @OneToMany(() => ChannelKeyword, (keyword) => keyword.channel)
  keywords: ChannelKeyword[];

  @OneToMany(() => ChannelModerationSetting, (moderationSettings) => moderationSettings.channel)
  moderationSettings: ChannelModerationSetting[];

  @OneToMany(() => ChannelPermit, (permit) => permit.channel)
  permits: ChannelPermit[];

  @OneToMany(() => ChannelTimer, (timer) => timer.channel)
  timers: ChannelTimer[];

  @OneToMany(() => UserStats, (stats) => stats.channel)
  usersStats: UserStats[];
}
