import { Column, Entity, JoinColumn, OneToMany, OneToOne, PrimaryColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { TwitchToken } from './TwitchToken.js';

export enum BotType {
  DEFAULT = 'default',
  custom = 'custom',
}

@Entity('bots')
export class Bot {
  @PrimaryColumn()
  id: string;

  @Column({
    type: 'enum',
    enum: BotType,
    default: BotType.DEFAULT,
  })
  type: BotType;

  @OneToOne(() => TwitchToken, t => t.bot)
  @JoinColumn()
  token: TwitchToken;

  @OneToMany(() => Channel, c => c.bot)
  channels: Channel[];
}