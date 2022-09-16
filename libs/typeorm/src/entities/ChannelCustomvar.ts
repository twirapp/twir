/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { Channel } from './Channel.js';

export enum CustomVarType {
  SCRIPT = 'SCRIPT',
  TEXT = 'TEXT',
}

@Index('channels_customvars_pkey', ['id'], { unique: true })
@Entity('channels_customvars', { schema: 'public' })
export class ChannelCustomvar {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'name' })
  name: string;

  @Column('text', { name: 'description', nullable: true })
  description: string | null;

  @Column('enum', { name: 'type', enum: CustomVarType })
  type: CustomVarType;

  @Column('text', { name: 'evalValue', nullable: true })
  evalValue: string | null;

  @Column('text', { name: 'response', nullable: true })
  response: string | null;

  @ManyToOne(() => Channel, (channels) => channels.customVar, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel: Channel;
}
