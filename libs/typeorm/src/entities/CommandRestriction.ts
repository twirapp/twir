import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryGeneratedColumn, type Relation } from 'typeorm';

import { type ChannelCommand } from './ChannelCommand.js';

export enum RestrictionType {
  WATCHED = 'WATCHED',
  MESSAGES = 'MESSAGES',
}

@Entity({ name: 'commands_restrictions' })
@Index('channels_restrictions_commandId_type', ['commandId', 'type'], { unique: true })
export class CommandRestriction {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  commandId: string;

  @ManyToOne('ChannelCommand', 'restrictions')
  @JoinColumn({ name: 'commandId' })
  command?: Relation<ChannelCommand>;

  @Column('enum', { enum: RestrictionType })
  type: RestrictionType;

  @Column('text')
  value: string;
}