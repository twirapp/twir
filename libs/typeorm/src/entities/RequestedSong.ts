import {
  Column,
  CreateDateColumn,
  DeleteDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn, PrimaryGeneratedColumn,
  Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

@Entity('channels_requested_songs')
export class RequestedSong {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('varchar')
  videoId: string;

  @Column('text')
  title: string;

  @Column('int4')
  duration: number;

  @CreateDateColumn()
  createdAt: Date;

  @Column()
  orderedById: string;

  @Column()
  orderedByName: string;

  @Column({ nullable: true })
  orderedByDisplayName: string | null;

  @Column()
  queuePosition: number;

  @ManyToOne('User', 'requestedSongs')
  @JoinColumn({ name: 'orderedById' })
  orderedBy?: Relation<User>;

  @Column()
  channelId: string;

  @ManyToOne('Channel', 'requestedSongs')
  @JoinColumn({ name: 'channelId' })
  channel?: Relation<Channel>;

  @DeleteDateColumn({ nullable: true })
  deletedAt: Date | null;
}
