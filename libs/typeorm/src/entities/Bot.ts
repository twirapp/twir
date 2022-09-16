/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, OneToMany, OneToOne } from 'typeorm';

import { Channel } from './Channel.js';
import { Token } from './Token.js';

export enum BotType {
  DEFAULT = 'DEFAULT',
  CUSTOM = 'CUSTOM',
}

@Index('bots_pkey', ['id'], { unique: true })
@Index('bots_tokenId_key', ['tokenId'], { unique: true })
@Entity('bots', { schema: 'public' })
export class Bot {
  @Column('text', { primary: true, name: 'id' })
  id: string;

  @Column('enum', { name: 'type', enum: BotType })
  type: BotType;

  @Column('text', { name: 'tokenId', nullable: true })
  tokenId: string | null;

  @OneToOne(() => Token, (tokens) => tokens.bots, {
    onDelete: 'SET NULL',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'tokenId', referencedColumnName: 'id' }])
  token: Token;

  @OneToMany(() => Channel, (channel) => channel.bot)
  channels: Channel[];
}
