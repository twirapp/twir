/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { ChannelCommand } from './ChannelCommand.js';

@Index('channels_commands_responses_pkey', ['id'], { unique: true })
@Entity('channels_commands_responses', { schema: 'public' })
export class CommandResponse {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'text', nullable: true })
  text: string | null;

  @ManyToOne(() => ChannelCommand, (command) => command.responses, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command: ChannelCommand;
}
