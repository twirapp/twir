/* eslint-disable import/no-cycle */
import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Entity('channels_messages')
export class ChannelChatMessage {
  @PrimaryColumn('text')
  messageId: string;

  @Column('text')
  @Index()
  channelId: string;

  @ManyToOne(() => Channel, _ => _.messages)
  @JoinColumn({ name: 'channelId' })
  channel?: Channel;

  @Column('text')
  @Index()
  userId: string;

  @ManyToOne(() => User, _ => _.messages)
  @JoinColumn({ name: 'userId' })
  user?: User;

  @Column('text')
  userName: string;

  @Column('text')
  text: string;

  @Column('boolean', { default: true })
  canBeDeleted: boolean;

  @CreateDateColumn()
  createdAt: Date;
}
