/* eslint-disable import/no-cycle */
import { Column, Entity, Index, JoinColumn, ManyToOne } from 'typeorm';

import { ChannelCommand } from './ChannelCommand.js';
import { User } from './User.js';

@Index('channels_commands_usages_pkey', ['id'], { unique: true })
@Entity('channels_commands_usages', { schema: 'public' })
export class CommandUsage {
  @Column('text', {
    primary: true,
    name: 'id',
    default: 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne(() => ChannelCommand, (command) => command.usages, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command: ChannelCommand;

  @ManyToOne(() => User, (user) => user.commandUsages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user: User;
}
