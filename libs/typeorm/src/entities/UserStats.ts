/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Index('users_stats_userId_channelId_key', ['channelId', 'userId'], {
  unique: true,
})
@Index('users_stats_pkey', ['id'], { unique: true })
@Entity('users_stats', { schema: 'public' })
export class UserStats {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'userId' })
  userId: string;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @Column('integer', { name: 'messages', default: 0 })
  messages: number;

  @Column('bigint', { name: 'watched', default: 0 })
  watched: bigint;

  @ManyToOne(() => Channel, (channel) => channel.usersStats, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;

  @ManyToOne(() => User, (user) => user.stats, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;
}
