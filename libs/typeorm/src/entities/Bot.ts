/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  OneToMany,
  OneToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { Token } from './Token';

export enum BotType {
  DEFAULT = 'DEFAULT',
  CUSTOM = 'CUSTOM',
}

@Entity('bots', { schema: 'public' })
export class Bot {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    unique: true,
    primaryKeyConstraintName: 'bots_pkey',
  })
  id: string;

  @Column('enum', { name: 'type', enum: BotType })
  type: BotType;

  @OneToOne(() => Token, _ => _.bot, {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([
    { name: 'tokenId', referencedColumnName: 'id', foreignKeyConstraintName: 'bots_tokenId_key' },
  ])
  token?: Token;

  @Index()
  @Column('text', { name: 'tokenId', nullable: true, unique: true })
  tokenId: string | null;

  @OneToMany(() => Channel, _ => _.bot)
  channels?: Channel[];
}
