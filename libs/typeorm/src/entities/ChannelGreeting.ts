/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';

@Index('channels_greetings_pkey', ['id'], { unique: true })
@Entity('channels_greetings', { schema: 'public' })
export class ChannelGreeting {
  @Column('text', {
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

  @ManyToOne(() => Channel, (channels) => channels.greetings, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;
}
