import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { CommandResponse } from './CommandResponse.js';
import { CommandUsage } from './CommandUsage.js';

export enum CooldownType {
  GLOBAL = 'global',
  PER_USER = 'user'
}

@Entity('channels_commands')
export class Command {
  @PrimaryGeneratedColumn('uuid')

  @Column()
  name: string;

  @Column({ default: 0 })
  cooldown: number;

  @Column({
    type: 'enum',
    enum: CooldownType,
    default: CooldownType.GLOBAL,
  })
  cooldownType: CooldownType;

  @Column({ default: true })
  enabled: boolean;

  @Column('simple-array', { default: '[]' })
  aliases: string[];

  @Column()
  description?: string;

  @Column({ default: true })
  visible: boolean;

  @ManyToOne(() => Channel, (c) => c.commands)
  @JoinColumn()
  channel: Channel;

  @ManyToOne(() => CommandResponse, (c) => c.command)
  @JoinColumn()
  responses: CommandResponse[];

  @OneToMany(() => CommandUsage, u => u.command)
  usages: CommandUsage[];
}