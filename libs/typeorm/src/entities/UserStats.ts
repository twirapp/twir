/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

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

  @ManyToOne('Channel', 'usersStats', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne('User', 'stats', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column('text', { name: 'userId' })
  userId: string;
}
