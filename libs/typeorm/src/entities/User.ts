/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToMany, OneToOne } from 'typeorm';

import { Channel } from './Channel.js';
import { ChannelPermit } from './ChannelPermit.js';
import { CommandUsage } from './CommandUsage.js';
import { DashboardAccess } from './DashboardAccess.js';
import { Notification } from './Notification.js';
import { Token } from './Token.js';
import { UserFile } from './UserFile.js';
import { UserStats } from './UserStats.js';
import { UserViewedNotification } from './UserViewedNotification.js';

@Index('users_pkey', ['id'], { unique: true })
@Index('users_tokenId_key', ['tokenId'], { unique: true })
@Entity('users', { schema: 'public' })
export class User {
  @Column('text', { primary: true, name: 'id' })
  id: string;

  @Column('text', { name: 'tokenId', nullable: true })
  tokenId: string | null;

  @Column('boolean', { name: 'isTester', default: false })
  isTester: boolean;

  @Column('boolean', { name: 'isBotAdmin', default: false })
  isBotAdmin: boolean;

  @OneToOne(() => Channel, (channel) => channel.user)
  channel: Channel;

  @OneToMany(() => CommandUsage, (commandUsages) => commandUsages.user)
  commandUsages: CommandUsage[];

  @OneToMany(() => DashboardAccess, (dashboardAccess) => dashboardAccess.user)
  dashboardAccess: DashboardAccess[];

  @OneToMany(() => ChannelPermit, (permit) => permit.user)
  permits: ChannelPermit[];

  @OneToMany(() => Notification, (notifications) => notifications.user)
  notifications: Notification[];

  @OneToOne(() => Token, (tokens) => tokens.users, { onDelete: 'SET NULL', onUpdate: 'CASCADE' })
  @JoinColumn([{ name: 'tokenId', referencedColumnName: 'id' }])
  token: Token;

  @OneToMany(() => UserFile, (file) => file.user)
  files: UserFile[];

  @OneToMany(() => UserStats, (stats) => stats.user)
  stats: UserStats[];

  @OneToMany(
    () => UserViewedNotification,
    (usersViewedNotifications) => usersViewedNotifications.user,
  )
  viewedNotifications: UserViewedNotification[];
}
