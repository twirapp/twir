/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Entity('channels_permits', { schema: 'public' })
export class ChannelPermit {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Channel, _ => _.permits, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId: string;

  @ManyToOne(() => User, _ => _.permits, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column()
  userId: string;
}
