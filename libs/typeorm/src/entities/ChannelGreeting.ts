/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
  PrimaryGeneratedColumn,
  Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';

@Entity('channels_greetings', { schema: 'public' })
export class ChannelGreeting {
  @PrimaryGeneratedColumn('uuid', {
    name: 'id',
  })
  id: string;

  @Column('text', { name: 'userId' })
  userId: string;

  @Column('boolean', { name: 'enabled', default: true })
  enabled: boolean;

  @Column('text', { name: 'text' })
  text: string;

  @ManyToOne('Channel', 'greetings', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;

  @Column('bool', { default: false })
  processed: boolean;
}
