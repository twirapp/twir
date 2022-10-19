import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn, Relation } from 'typeorm';

import { ChannelTimer } from './ChannelTimer.js';

@Entity('channels_timers_responses', { schema: 'public' })
export class ChannelTimerResponse {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('text')
  text: string;

  @Column('bool', { default: true })
  isAnnounce: boolean;

  @ManyToOne('ChannelTimer', 'responses')
  @JoinColumn({ name: 'timerId' })
  timer?: Relation<ChannelTimer>;

  @Column()
  timerId: string;
}
