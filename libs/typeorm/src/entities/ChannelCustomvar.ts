/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

import { Channel } from './Channel';

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

  @Column('text', { name: 'evalValue', nullable: false, default: '' })
  evalValue: string;

  @Column('text', { name: 'response', nullable: false, default: '' })
  response: string;

  @ManyToOne(() => Channel, _ => _.customVar, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column()
  channelId?: string;
}
