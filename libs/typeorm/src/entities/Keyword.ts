import { Column, Entity, JoinColumn, ManyToMany, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel.js';

@Entity('channels_keywords')
export class Keyword {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToMany(() => Channel, c => c.keywords)
  @JoinColumn()
  channel: Channel;

  @Column()
  text: string;

  @Column()
  response: string;

  @Column({ default: true })
  enabled: boolean;

  @Column({ default: 0 })
  cooldown: number;
}