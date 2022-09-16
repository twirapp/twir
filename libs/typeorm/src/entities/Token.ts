/* eslint-disable import/no-cycle */
import { Column, Entity, Index, OneToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Bot } from './Bot.js';
import { type User } from './User.js';

@Entity('tokens', { schema: 'public' })
export class Token {
  @PrimaryColumn('text', {
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

  @OneToOne('Bot', 'token')
  bots: Relation<Bot>;

  @OneToOne('User', 'token')
  users: Relation<User>;
}
