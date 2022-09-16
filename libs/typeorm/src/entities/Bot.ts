/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToMany, OneToOne, PrimaryColumn, Relation, Unique } from 'typeorm';

import { type Channel } from './Channel.js';
import { type Token } from './Token.js';

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
  
  @Index()
  @Column('text', { name: 'tokenId', nullable: true, unique: true })
  tokenId: string | null;

  @OneToOne('Token', 'bots', {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'tokenId', referencedColumnName: 'id', foreignKeyConstraintName: 'bots_tokenId_key' }])
  token: Relation<Token>;

  @OneToMany('Channel', 'bot')
  channels: Relation<Channel[]>;
}
