/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';

@Index('channels_timers_pkey', ['id'], { unique: true })
@Entity('channels_timers', { schema: 'public' })
export class ChannelTimer {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('character varying', { name: 'name', length: 255 })
  name: string;

  @Column('boolean', { name: 'enabled', default: false })
  enabled: boolean;

  @Column('jsonb', { name: 'responses', default: [] })
  responses: string[];

  @Column('integer', { name: 'last', default: 0 })
  last: number;

  @Column('integer', { name: 'timeInterval', default: 0 })
  timeInterval: number;

  @Column('integer', { name: 'messageInterval', default: 0 })
  messageInterval: number;

  @Column('integer', { name: 'lastTriggerMessageNumber', default: 0 })
  lastTriggerMessageNumber: number;

  @ManyToOne(() => Channel, (channel) => channel.timers, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;
}
