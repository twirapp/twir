import { Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Entity('channels_permits')
export class Permit {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Channel, c => c.permits)
  @JoinColumn()
  channel: Channel;

  @ManyToOne(() => User, u => u.permits)
  @JoinColumn()
  user: User;
}