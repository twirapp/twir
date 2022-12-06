/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Entity('users_online', { schema: 'public' })
export class UserOnline {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @ManyToOne(() => Channel, _ => _.onlineUsers, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne(() => User, _ => _.online, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column('text', { nullable: true })
  userId: string | null;

  @Column('text', { nullable: true })
  userName: string | null;
}
