/* eslint-disable import/no-cycle */
import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  OneToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { ChannelEvent } from '../ChannelEvent';

@Entity({
  name: 'channel_events_donations',
})
export class ChannelDonationEvent {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @OneToOne(() => ChannelEvent, _ => _.donation)
  @JoinColumn({ name: 'eventId' })
  event?: ChannelEvent;

  @Column()
  eventId: string;

  @Column({ nullable: true })
  fromUserId: string | null;

  @Column({ nullable: true })
  toUserId: string | null;

  @Column('numeric')
  amount: number;

  @Column()
  currency: string;

  @Column({ nullable: true })
  username: string | null;

  @Column({ nullable: true })
  message: string | null;

  @CreateDateColumn()
  createdAt: Date;


  @Column({ nullable: true, unique: true })
  donateId: string | null;
}
