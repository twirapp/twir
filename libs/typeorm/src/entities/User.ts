/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToMany, OneToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { type ChannelPermit } from './ChannelPermit.js';
import { type CommandUsage } from './CommandUsage.js';
import { type DashboardAccess } from './DashboardAccess.js';
import { type Notification } from './Notification.js';
import { type Token } from './Token.js';
import { type UserFile } from './UserFile.js';
import { type UserStats } from './UserStats.js';
import { type UserViewedNotification } from './UserViewedNotification.js';

@Entity('users', { schema: 'public' })
export class User {
  @PrimaryColumn('text', { primary: true, name: 'id' })
  id: string;

  @Index()
  @Column('text', { name: 'tokenId', nullable: true })
  tokenId: string | null;

  @Column('boolean', { name: 'isTester', default: false })
  isTester: boolean;

  @Column('boolean', { name: 'isBotAdmin', default: false })
  isBotAdmin: boolean;

  @OneToOne('Channel', 'user')
  channel: Relation<Channel>;

  @OneToMany('CommandUsage', 'user')
  commandUsages: Relation<CommandUsage[]>;

  @OneToMany('DashboardAccess', 'user')
  dashboardAccess: Relation<DashboardAccess[]>;

  @OneToMany('ChannelPermit', 'user')
  permits: Relation<ChannelPermit[]>;

  @OneToMany('Notification', 'user')
  notifications: Relation<Notification[]>;

  @OneToOne('Token', 'users', { onDelete: 'SET NULL', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'tokenId', referencedColumnName: 'id' }])
  token: Relation<Token>;

  @OneToMany('UserFile', 'user')
  files: Relation<UserFile[]>;

  @OneToMany('UserStats', 'user')
  stats: Relation<UserStats[]>;

  @OneToMany('UserViewedNotification', 'user')
  viewedNotifications: Relation<UserViewedNotification[]>;
}
