import { Column, Entity, OneToOne, PrimaryGeneratedColumn, type Relation } from 'typeorm';

import { type Bot } from './Bot.js';
import { type User } from './User.js';

@Entity('tokens', { schema: 'public' })
export class Token {
  @PrimaryGeneratedColumn('uuid', {
    name: 'id',
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
  bot?: Relation<Bot>;

  @OneToOne('User', 'token')
  user?: Relation<User>;
}
