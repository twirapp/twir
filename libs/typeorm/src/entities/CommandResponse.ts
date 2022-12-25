/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { ChannelCommand } from './ChannelCommand';

@Entity('channels_commands_responses', { schema: 'public' })
export class CommandResponse {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'text', nullable: true })
  text: string | null;

  @ManyToOne(() => ChannelCommand, _ => _.responses, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command?: ChannelCommand;

  @Column()
  commandId: string;

  @Column('int', { default: 0 })
  order: number;
}
