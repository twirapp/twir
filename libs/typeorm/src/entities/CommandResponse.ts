/* eslint-disable import/no-cycle */
import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, Relation } from 'typeorm';

import { type ChannelCommand } from './ChannelCommand.js';

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

  @ManyToOne('ChannelCommand', 'responses', {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command?: Relation<ChannelCommand>;

  @Column()
  commandId: string;
}
