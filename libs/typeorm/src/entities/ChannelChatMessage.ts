import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
  type Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

@Entity('channels_messages')
export class ChannelChatMessage {
  @PrimaryColumn('text')
  messageId: string;

  @Column('text')
  @Index()
  channelId: string;

  @ManyToOne('Channel', 'messages')
  @JoinColumn({ name: 'channelId' })
  channel?: Relation<Channel>;

  @Column('text')
  @Index()
  userId: string;

  @ManyToOne('User', 'messages')
  @JoinColumn({ name: 'userId' })
  user?: Relation<User>;

  @Column('text')
  userName: string;

  @Column('text')
  text: string;

  @Column('boolean', { default: true })
  canBeDeleted: boolean;

  @CreateDateColumn()
  createdAt: Date;
}
