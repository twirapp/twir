/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  OneToMany,
  OneToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { ChannelChatMessage } from './ChannelChatMessage';
import { ChannelPermit } from './ChannelPermit';
import { CommandUsage } from './CommandUsage';
import { DashboardAccess } from './DashboardAccess';
import { Notification } from './Notification';
import { Token } from './Token';
import { UserFile } from './UserFile';
import { UserOnline } from './UserOnline';
import { UserStats } from './UserStats';
import { UserViewedNotification } from './UserViewedNotification';

@Entity('users', { schema: 'public' })
export class User {
  @PrimaryColumn('text', { primary: true, name: 'id' })
  id: string;

  @Column('boolean', { name: 'isTester', default: false })
  isTester: boolean;

  @Column('boolean', { name: 'isBotAdmin', default: false })
  isBotAdmin: boolean;

  @OneToOne(() => Channel, _ => _.user)
  channel?: Channel;

  @OneToMany(() => CommandUsage, _ => _.user)
  commandUsages?: CommandUsage[];

  @OneToMany(() => DashboardAccess, _ => _.user)
  dashboardAccess?: DashboardAccess[];

  @OneToMany(() => ChannelPermit, _ => _.user)
  permits?: ChannelPermit[];

  @OneToMany(() => Notification, _ => _.user)
  notifications?: Notification[];

  @OneToOne(() => Token, _ => _.user, {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'tokenId', referencedColumnName: 'id' }])
  token?: Token;

  @Index()
  @Column('text', { name: 'tokenId', nullable: true })
  tokenId: string | null;

  @Column('uuid', { generated: 'uuid' })
  apiKey: string;

  @OneToMany(() => UserFile, _ => _.user)
  files?: UserFile[];

  @OneToMany(() => UserStats, _ => _.user)
  stats?: UserStats[];

  @OneToOne(() => UserOnline, _ => _.user)
  online?: UserOnline;

  @OneToMany(() => UserViewedNotification, _ => _.user)
  viewedNotifications?: UserViewedNotification[];

  @OneToMany(() => ChannelChatMessage, _ => _.user)
  messages?: ChannelChatMessage[];
}
