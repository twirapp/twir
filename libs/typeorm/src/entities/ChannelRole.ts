/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel';
import { ChannelRoleUser } from './ChannelRoleUser';

export enum RoleType {
  BROADCASTER = 'BROADCASTER',
  MODERATOR = 'MODERATOR',
  SUBSCRIBER = 'SUBSCRIBER',
  VIP = 'VIP',
  CUSTOM = 'CUSTOM',
}

export enum RoleFlags {
  CAN_ACCESS_DASHBOARD = 'CAN_ACCESS_DASHBOARD',

  UPDATE_CHANNEL_TITLE = 'UPDATE_CHANNEL_TITLE',
  UPDATE_CHANNEL_CATEGORY = 'UPDATE_CHANNEL_CATEGORY',

  VIEW_COMMANDS = 'VIEW_COMMANDS',
  MANAGE_COMMANDS = 'MANAGE_COMMANDS',

  VIEW_KEYWORDS = 'VIEW_KEYWORDS',
  MANAGE_KEYWORDS = 'MANAGE_KEYWORDS',

  VIEW_TIMERS = 'VIEW_TIMERS',
  MANAGE_TIMERS = 'MANAGE_TIMERS',

  VIEW_INTEGRATIONS = 'VIEW_INTEGRATIONS',
  MANAGE_INTEGRATIONS = 'MANAGE_INTEGRATIONS',

  VIEW_SONG_REQUESTS = 'VIEW_SONG_REQUESTS',
  MANAGE_SONG_REQUESTS = 'MANAGE_SONG_REQUESTS',

  VIEW_MODERATION = 'VIEW_MODERATION',
  MANAGE_MODERATION = 'MANAGE_MODERATION',

  VIEW_VARIABLES = 'VIEW_VARIABLES',
  MANAGE_VARIABLES = 'MANAGE_VARIABLES',

  VIEW_GREETINGS = 'VIEW_GREETINGS',
  MANAGE_GREETINGS = 'MANAGE_GREETINGS',
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

  @Column('enum', { enum: RoleType, default: RoleType.CUSTOM })
  type: RoleType;

  @Column({
    type: 'enum',
    enum: RoleFlags,
    array: true,
    default: [],
  })
  permissions: RoleFlags[];

  @OneToMany(() => ChannelRoleUser, _ => _.role)
  users?: ChannelRoleUser[];
}