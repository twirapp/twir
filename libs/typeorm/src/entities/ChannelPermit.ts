import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  type Relation,
} from 'typeorm';

import { type Channel } from './Channel.js';
import { type User } from './User.js';

@Entity('channels_permits', { schema: 'public' })
export class ChannelPermit {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne('Channel', 'permits', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;

  @ManyToOne('User', 'permits', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column()
  userId: string;
}
