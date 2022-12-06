/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';

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

  @ManyToOne(() => Channel, _ => _.greetings, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId: string;

  @Column('bool', { default: false })
  processed: boolean;

  @Column('bool', { default: true })
  isReply: boolean;
}
