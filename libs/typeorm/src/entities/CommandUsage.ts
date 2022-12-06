/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryColumn,
} from 'typeorm';

import { ChannelCommand } from './ChannelCommand';
import { User } from './User';

@Entity('channels_commands_usages', { schema: 'public' })
export class CommandUsage {
  @PrimaryColumn('text', {
    primary: true,
    name: 'id',
    default: () => 'gen_random_uuid()',
  })
  id: string;

  @Column('text', { name: 'channelId' })
  channelId: string;

  @ManyToOne(() => ChannelCommand, _ => _.usages, {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command?: ChannelCommand;

  @Column()
  commandId: string;

  @ManyToOne(() => User, _ => _.commandUsages, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: User;

  @Column()
  userId: string;
}
