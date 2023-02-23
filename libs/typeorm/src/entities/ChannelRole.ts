/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel';
import { ChannelRolePermission } from './ChannelRolePermission';
import { ChannelRoleUser } from './ChannelRoleUser';

export enum ChannelRoleType {
  BROADCASTER = 'BROADCASTER',
  MODERATOR = 'MODERATOR',
  SUBSCRIBER = 'SUBSCRIBER',
  VIP = 'VIP',
  // FOLLOWER = 'FOLLOWER',
  // VIEWER = 'VIEWER',
  CUSTOM = 'CUSTOM',
}

@Entity('channels_roles')
export class ChannelRole {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Channel, channel => channel.roles)
  @JoinColumn({ name: 'channelId' })
  channel?: Channel;

  @Column()
  channelId: string;

  @Column()
  name: string;

  @Column('enum', { enum: ChannelRoleType, default: ChannelRoleType.CUSTOM })
  type: ChannelRoleType;

  @Column({ default: false })
  system: boolean;

  @OneToMany(() => ChannelRolePermission, _ => _.role)
  permissions: ChannelRolePermission[];

  @OneToMany(() => ChannelRoleUser, _ => _.role)
  users?: ChannelRoleUser[];
}