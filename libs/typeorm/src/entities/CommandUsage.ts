import { Column, Entity, JoinColumn, ManyToOne, PrimaryColumn, type Relation } from 'typeorm';

import { type ChannelCommand } from './ChannelCommand.js';
import { type User } from './User.js';

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

  @ManyToOne('ChannelCommand', 'usages', {
    onDelete: 'CASCADE',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'commandId', referencedColumnName: 'id' }])
  command?: Relation<ChannelCommand>;

  @Column()
  commandId: string;

  @ManyToOne('User', 'commandUsages', {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'userId', referencedColumnName: 'id' }])
  user?: Relation<User>;

  @Column()
  userId: string;
}
