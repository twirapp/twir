import { Column, Entity, JoinColumn, OneToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Bot } from './Bot.js';
import { User } from './User.js';

@Entity('twitch_tokens')
export class TwitchToken {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  accessToken: string;

  @Column()
  refreshToken: string;

  @Column()
  expiresInt: number;

  @Column()
  obtainmentTimestamp: Date;

  @OneToOne(() => Bot, b => b.token)
  bot?: Bot;

  @OneToOne(() => User, u => u.token)
  user?: User;
}