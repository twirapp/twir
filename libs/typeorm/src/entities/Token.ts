/* eslint-disable import/no-cycle */
import { Column, Entity, Index, OneToOne } from 'typeorm';

import { Bot } from './Bot.js';
import { User } from './User.js';

@Index('tokens_pkey', ['id'], { unique: true })
@Entity('tokens', { schema: 'public' })
export class Token {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'accessToken' })
  accessToken: string;

  @Column('text', { name: 'refreshToken' })
  refreshToken: string;

  @Column('integer', { name: 'expiresIn' })
  expiresIn: number;

  @Column('timestamp without time zone', { name: 'obtainmentTimestamp' })
  obtainmentTimestamp: Date;

  @OneToOne(() => Bot, (bot) => bot.token)
  bots: Bot;

  @OneToOne(() => User, (user) => user.token)
  users: User;
}
