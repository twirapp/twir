import { Column, Entity, JoinColumn, OneToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Entity('users_stats')
export class UserStats {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  messags: number;

  @Column('bigint')
  watched: string;

  @OneToOne(() => User, (u) => u.stats)
  @JoinColumn()
  user: User;

  @OneToOne(() => Channel, (u) => u.usersStats)
  @JoinColumn()
  channel: Channel;
}