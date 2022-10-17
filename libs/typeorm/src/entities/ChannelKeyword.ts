/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
  PrimaryGeneratedColumn,
  Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';

@Index('channels_keywords_channelId_text_key', ['channelId', 'text'], {
  unique: true,
})
@Entity('channels_keywords', { schema: 'public' })
export class ChannelKeyword {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('text', { name: 'text' })
  text: string;

  @Column('text', { name: 'response' })
  response: string;

  @Column('boolean', { name: 'enabled', default: true })
  enabled: boolean;

  @Column('integer', { name: 'cooldown', nullable: true, default: 0 })
  cooldown: number | null;

  @ManyToOne('Channel', 'keywords', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @Column('timestamp', { nullable: true })
  cooldownExpireAt: Date | null;
}
