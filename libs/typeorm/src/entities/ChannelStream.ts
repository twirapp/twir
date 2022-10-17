import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import type { Channel } from './Channel.js';

@Entity('channels_streams')
export class ChannelStream {
  @PrimaryColumn()
  id: string;

  @Column('text')
  userId: string;

  @Column('text')
  userLogin: string;

  @Column('text')
  userName: string;

  @Column({ type: 'int' })
  gameId: string;

  @Column('text')
  gameName: string;

  @Column('text', { array: true, default: [], nullable: true })
  communityIds: string[] | null;

  @Column({ type: 'text' })
  type: 'live' | 'vodcast' | '';

  @Column('text')
  title: string;

  @Column({ type: 'int' })
  viewerCount: number;

  @Column('timestamp')
  startedAt: Date;

  @Column('text')
  language: string;

  @Column('text')
  thumbnailUrl: string;

  @Column('text', { nullable: true, array: true, default: [] })
  tagIds: string[] | null;

  @Column('bool')
  isMature: boolean;

  @Column({ default: 0, type: 'int' })
  parsedMessages: number;

  @ManyToOne('Channel', 'streams')
  @JoinColumn({ name: 'userId' })
  channel?: Relation<Channel>;
}
