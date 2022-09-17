/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryColumn, PrimaryGeneratedColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { type CommandResponse } from './CommandResponse.js';
import { type CommandUsage } from './CommandUsage.js';

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

@Index('channels_commands_name_channelId_key', ['channelId', 'name'], { unique: true })
@Entity('channels_commands', { schema: 'public' })
export class ChannelCommand {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Index()
  @Column('text', { name: 'name' })
  name: string;

  @Column('integer', { name: 'cooldown', nullable: true, default: 0 })
  cooldown: number | null;

  @Column('enum', { name: 'cooldownType', enum: CooldownType, default: CooldownType.GLOBAL })
  cooldownType: CooldownType;

  @Column('boolean', { name: 'enabled', default: true })
  enabled: boolean;

  @Column('text', { name: 'aliases', array: true, default: [] })
  aliases: string[];

  @Column('text', { name: 'description', nullable: true })
  description: string | null;

  @Column('boolean', { name: 'visible', default: true })
  visible: boolean;

  
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
  
  @ManyToOne('Channel', 'commands', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;
  
  @Index()
  @Column('text', { name: 'channelId' })
  channelId: string;
  
  @OneToMany('CommandResponse', 'command')
  responses?: Relation<CommandResponse[]>;
  
  @OneToMany('CommandUsage', 'command')
  usages?: Relation<CommandUsage[]>;
}
