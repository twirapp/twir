/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { ChannelTimerResponse } from './ChannelTimerResponse';

@Entity('channels_timers', { schema: 'public' })
export class ChannelTimer {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('character varying', { name: 'name', length: 255 })
  name: string;

  @Column('boolean', { name: 'enabled', default: false })
  enabled: boolean;

  @OneToMany(() => ChannelTimerResponse, _ => _.timer, {
    cascade: true,
  })
  responses: ChannelTimerResponse[];

  @Column('integer', { name: 'timeInterval', default: 0 })
  timeInterval: number;

  @Column('integer', { name: 'messageInterval', default: 0 })
  messageInterval: number;

  @Column('integer', { name: 'lastTriggerMessageNumber', default: 0 })
  lastTriggerMessageNumber: number;

  @ManyToOne(() => Channel, _ => _.timers, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId: string;
}
