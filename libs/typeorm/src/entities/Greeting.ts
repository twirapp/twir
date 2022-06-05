import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';
import { User } from './User.js';

@Entity('channels_greetings')
export class Greeting {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ default: true })
  enabled: boolean;

  @Column()
  text: string;

  @ManyToOne(() => User, u => u.greetings)
  @JoinColumn()
  user: User;

  @ManyToOne(() => Channel, c => c.greetings)
  @JoinColumn()
  channel: Channel;
}