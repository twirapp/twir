/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { ChannelTimerResponse } from './ChannelTimerResponse.js';

@Entity('channels_timers', { schema: 'public' })
export class ChannelTimer {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('character varying', { name: 'name', length: 255 })
  name: string;

  @Column('boolean', { name: 'enabled', default: false })
  enabled: boolean;

  @OneToMany('ChannelTimerResponse', 'timer', {
    cascade: true,
  })
  responses: Relation<ChannelTimerResponse>;

  @Column('integer', { name: 'timeInterval', default: 0 })
  timeInterval: number;

  @Column('integer', { name: 'messageInterval', default: 0 })
  messageInterval: number;

  @Column('integer', { name: 'lastTriggerMessageNumber', default: 0 })
  lastTriggerMessageNumber: number;

  @ManyToOne('Channel', 'timers', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;
}
