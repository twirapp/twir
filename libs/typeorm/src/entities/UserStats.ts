/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Index('users_stats_userId_channelId_key', ['channelId', 'userId'], {
  unique: true,
})
@Entity('users_stats', { schema: 'public' })
export class UserStats {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('integer', { name: 'messages', default: 0 })
  messages: number;

  @Column('bigint', { name: 'watched', default: 0 })
  watched: bigint;

  @Column('bigint', { default: 0 })
  usedChannelPoints: bigint;

  @ManyToOne(() => Channel, _ => _.usersStats, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne(() => User, _ => _.stats, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column('text', { name: 'userId' })
  userId: string;
}
