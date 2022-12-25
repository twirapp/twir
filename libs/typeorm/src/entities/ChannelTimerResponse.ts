/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { ChannelTimer } from './ChannelTimer';

@Entity('channels_timers_responses', { schema: 'public' })
export class ChannelTimerResponse {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('text')
  text: string;

  @Column('bool', { default: true })
  isAnnounce: boolean;

  @ManyToOne(() => ChannelTimer, _ => _.responses, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'timerId' })
  timer?: ChannelTimer;

  @Column()
  timerId: string;
}
