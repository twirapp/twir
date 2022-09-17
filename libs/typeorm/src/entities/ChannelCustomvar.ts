/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm';

import { type Channel } from './Channel.js';

export enum CustomVarType {
  SCRIPT = 'SCRIPT',
  TEXT = 'TEXT',
}

@Entity('channels_customvars', { schema: 'public' })
export class ChannelCustomvar {
  @PrimaryGeneratedColumn('uuid', {
    name: 'id',
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

  @ManyToOne('Channel', 'customVar', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId?: string;
}
