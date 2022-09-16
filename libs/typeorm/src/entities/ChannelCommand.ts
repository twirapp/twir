/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany } from 'typeorm';

import { Channel } from './Channel.js';
import { CommandResponse } from './CommandResponse.js';
import { CommandUsage } from './CommandUsage.js';

export enum CooldownType {
  GLOBAL = 'GLOBAL',
  PER_USER = 'PER_USER',
}

export enum CommandPermission {
  BROADCASTER = 'BROADCASTER',
  MODERATOR = 'MODERATOR',
  SUBSCRIBER = 'SUBSCRIBER',
  VIP = 'VIP',
  VIEWER = 'VIEWER',
  FOLLOWER = 'FOLLOWER',
}

export enum CommandModule {
  CUSTOM = 'CUSTOM',
  DOTA = 'DOTA',
  CHANNEL = 'CHANNEL',
  MODERATION = 'MODERATION',
}

@Index('channels_commands_channelId_idx', ['channelId'], {})
@Index('channels_commands_name_channelId_key', ['channelId', 'name'], { unique: true })
@Index('channels_commands_pkey', ['id'], { unique: true })
@Index('channels_commands_name_idx', ['name'], {})
@Entity('channels_commands', { schema: 'public' })
export class ChannelCommand {
  @Column('text', { primary: true, name: 'id', default: () => 'gen_random_uuid()' })
  id: string;

  @Column('text', { name: 'name' })
  name: string;

  @Column('integer', { name: 'cooldown', nullable: true, default: 0 })
  cooldown: number | null;

  @Column('enum', { name: 'cooldownType', enum: CooldownType, default: CooldownType.GLOBAL })
  cooldownType: CooldownType;

  @Column('boolean', { name: 'enabled', default: true })
  enabled: boolean;

  @Column('jsonb', { name: 'aliases', nullable: true, default: [] })
  aliases: object | null;

  @Column('text', { name: 'description', nullable: true })
  description: string | null;

  @Column('boolean', { name: 'visible', default: true })
  visible: boolean;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @Column('enum', {
    name: 'permission',
    enum: CommandPermission,
  })
  permission: CommandPermission;

  @Column('boolean', { name: 'default', default: false })
  default: boolean;

  @Column('text', { name: 'defaultName', nullable: true })
  defaultName: string | null;

  @Column('enum', {
    name: 'module',
    enum: CommandModule,
    default: CommandModule.CUSTOM,
  })
  module: CommandModule;

  @ManyToOne(() => Channel, (channel) => channel.commands, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;

  @OneToMany(() => CommandResponse, (response) => response.command)
  responses: CommandResponse[];

  @OneToMany(() => CommandUsage, (usage) => usage.command)
  usages: CommandUsage[];
}
