/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToMany, OneToOne, PrimaryColumn } from 'typeorm';

import { Channel } from './Channel';
import { ChannelChatMessage } from './ChannelChatMessage';
import { ChannelEmoteUsage } from './ChannelEmoteUsage';
import { ChannelGiveawayParticipant } from './ChannelGiveawayParticipant';
import { ChannelModuleSettings } from './ChannelModuleSettings';
import { ChannelPermit } from './ChannelPermit';
import { ChannelRoleUser } from './ChannelRoleUser';
import { CommandUsage } from './CommandUsage';
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

  @OneToOne(() => Channel, (_) => _.user)
  channel?: Channel;

  @OneToMany(() => CommandUsage, (_) => _.user)
  commandUsages?: CommandUsage[];

  @OneToMany(() => ChannelPermit, (_) => _.user)
  permits?: ChannelPermit[];

  @OneToMany(() => Notification, (_) => _.user)
  notifications?: Notification[];

  @OneToOne(() => Token, (_) => _.user, {
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

  @OneToMany(() => UserFile, (_) => _.user)
  files?: UserFile[];

  @OneToMany(() => UserStats, (_) => _.user)
  stats?: UserStats[];

  @OneToOne(() => UserOnline, (_) => _.user)
  online?: UserOnline;

  @OneToMany(() => UserViewedNotification, (_) => _.user)
  viewedNotifications?: UserViewedNotification[];

  @OneToMany(() => ChannelChatMessage, (_) => _.user)
  messages?: ChannelChatMessage[];

  @OneToMany(() => ChannelEmoteUsage, (_) => _.channel)
  emotesUsages?: ChannelEmoteUsage[];

  @OneToOne(() => ChannelModuleSettings, (_) => _.user)
  ttsSettings?: ChannelModuleSettings;

  @OneToMany(() => ChannelRoleUser, (_) => _.user)
  channelRoleUsers?: ChannelRoleUser[];

  @OneToMany(() => ChannelGiveawayParticipant, (_) => _.user)
  giveaway_participants?: ChannelGiveawayParticipant[];
}
