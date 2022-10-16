/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

@Entity('users_online', { schema: 'public' })
export class UserOnline {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @ManyToOne('Channel', 'onlineUsers', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne('User', 'online', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column('text', { nullable: true })
  userId: string | null;

  @Column('text', { nullable: true })
  userName: string | null;
}
