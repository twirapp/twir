import {
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  OneToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { ChannelEvent } from '../ChannelEvent.js';

@Entity({
  name: 'channel_events_follows',
})
export class ChannelFollowEvent {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @OneToOne('ChannelEvent', 'follow')
  @JoinColumn({ name: 'eventId' })
  event?: ChannelEvent;

  @Column()
  eventId: string;

  @Column()
  fromUserId: string;

  @Column()
  toUserId: string;

  @CreateDateColumn()
  createdAt: Date;
}
