import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';

@Entity('channels_timers')
export class Timer {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Channel, c => c.timers)
  @JoinColumn()
  channel: Channel;

  @Column()
  name: string;

  @Column({ default: true })
  enabled: boolean;

  @Column('simple-array', { default: '[]' })
  responses: string[];

  @Column({ default: 0 })
  last: number;

  @Column({ default: 0 })
  timeInterval: number;

  @Column({ default: 0 })
  messageInterval: number;
}