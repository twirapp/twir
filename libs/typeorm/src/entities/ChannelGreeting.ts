/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';

@Entity('channels_greetings', { schema: 'public' })
export class ChannelGreeting {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'userId' })
  userId: string;

  @Column('boolean', { name: 'enabled', default: true })
  enabled: boolean;

  @Column('text', { name: 'text' })
  text: string;

  @ManyToOne('Channel', 'greetings', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Relation<Channel>;
}
