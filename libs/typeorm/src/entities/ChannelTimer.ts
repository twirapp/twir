/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn, Relation } from 'typeorm';

import { type Channel } from './Channel.js';

@Entity('channels_timers', { schema: 'public' })
export class ChannelTimer {
  @PrimaryGeneratedColumn('uuid')
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

  @ManyToOne('Channel', 'timers', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Relation<Channel>;

  @Column()
  channelId: string;
}
