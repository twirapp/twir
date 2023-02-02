/* eslint-disable import/no-cycle */
import { Column, CreateDateColumn, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel';
import { User } from './User';

@Entity('channels_emotes_usages', { schema: 'public' })
export class ChannelEmoteUsage {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @ManyToOne(() => Channel, _ => _.emotesUsages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId' }])
  channel?: Channel;

  @Column()
  channelId?: string;

  @ManyToOne(() => User, _ => _.emotesUsages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId' }])
  user?: User;

  @Column()
  userId: string;

  @CreateDateColumn()
  createdAt: Date;

  @Column()
  emote: string;
}
