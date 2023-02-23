/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  OneToOne,
  PrimaryColumn,
} from 'typeorm';

import { Bot } from './Bot';
import { ChannelChatMessage } from './ChannelChatMessage';
import { ChannelCommand } from './ChannelCommand';
import { ChannelCommandGroup } from './ChannelCommandGroup';
import { ChannelCustomvar } from './ChannelCustomvar';
import { ChannelDotaAccount } from './ChannelDotaAccount';
import { ChannelEmoteUsage } from './ChannelEmoteUsage';
import { ChannelEvent } from './ChannelEvent';
import { ChannelGreeting } from './ChannelGreeting';
import { ChannelInfoHistory } from './ChannelInfoHistory';
import { ChannelIntegration } from './ChannelIntegration';
import { ChannelKeyword } from './ChannelKeyword';
import { ChannelModerationSetting } from './ChannelModerationSetting';
import { ChannelPermit } from './ChannelPermit';
import { ChannelRole } from './ChannelRole';
import { ChannelRoleUser } from './ChannelRoleUser';
import { ChannelStream } from './ChannelStream';
import { ChannelTimer } from './ChannelTimer';
import { DashboardAccess } from './DashboardAccess';
import { Event } from './events/Event';
import { User } from './User';
import { UserOnline } from './UserOnline';
import { UserStats } from './UserStats';

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

  @ManyToOne(() => Bot, _ => _.channels, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'botId', referencedColumnName: 'id' }])
  bot?: Bot;

  @Column()
  botId: string;

  @Column('bool', { default: false })
  isBotMod: boolean;

  @OneToOne(() => User, _ => _.channel, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'id', referencedColumnName: 'id' }])
  user?: User;

  @OneToMany(() => ChannelCommand, _ => _.channel)
  commands?: ChannelCommand[];

  @OneToMany(() => ChannelCustomvar, _ => _.channel)
  customVar?: ChannelCustomvar[];

  @OneToMany(() => DashboardAccess, _ => _.channel)
  dashboardAccess?: DashboardAccess[];

  @OneToMany(() => ChannelDotaAccount, _ => _.channel)
  dotaAccounts?: ChannelDotaAccount[];

  @OneToMany(() => ChannelGreeting, _ => _.channel)
  greetings?: ChannelGreeting[];

  @OneToMany(() => ChannelIntegration, _ => _.channel)
  channelsIntegrations?: ChannelIntegration[];

  @OneToMany(() => ChannelKeyword, _ => _.channel)
  keywords?: ChannelKeyword[];

  @OneToMany(() => ChannelModerationSetting, _ => _.channel)
  moderationSettings?: ChannelModerationSetting[];

  @OneToMany(() => ChannelPermit, _ => _.channel)
  permits?: ChannelPermit[];

  @OneToMany(() => ChannelTimer, _ => _.channel)
  timers?: ChannelTimer[];

  @OneToMany(() => UserStats, _ => _.channel)
  usersStats?: UserStats[];

  @OneToMany(() => ChannelEvent, _ => _.channel)
  eventsList?: ChannelEvent[];

  @OneToMany(() => Event, _ => _.channel)
  events: Event[];

  @OneToMany(() => ChannelStream, _ => _.channel)
  streams?: ChannelStream[];

  @OneToMany(() => UserOnline, _ => _.channel)
  onlineUsers?: UserOnline[];

  @OneToMany(() => ChannelChatMessage, _ => _.channel)
  messages?: ChannelChatMessage[];

  @OneToMany(() => ChannelEmoteUsage, _ => _.channel)
  emotesUsages?: ChannelEmoteUsage[];

  @OneToMany(() => ChannelInfoHistory, _ => _.channel)
  infoHistories?: ChannelInfoHistory[];

  @OneToMany(() => ChannelCommandGroup, _ => _.channel)
  commandsGroups?: ChannelCommandGroup[];

  @OneToMany(() => ChannelRole, _ => _.channel)
  roles?: ChannelRole[];
}
