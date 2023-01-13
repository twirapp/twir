/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  OneToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Bot } from './Bot';
import { User } from './User';

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

  @Column('text', { array: true, default: [], nullable: true })
  scopes: string[];

  @OneToOne(() => Bot, _ => _.token)
  bot?: Bot;

  @OneToOne(() => User, _ => _.token)
  user?: User;
}
