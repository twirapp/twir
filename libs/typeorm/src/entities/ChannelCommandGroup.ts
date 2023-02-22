import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Channel } from './Channel';
// eslint-disable-next-line import/no-cycle
import { ChannelCommand } from './ChannelCommand';

@Entity({ name: 'channels_commands_groups' })
export class ChannelCommandGroup {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  channelId: string;

  @ManyToOne(() => Channel, _ => _.commandsGroups)
  @JoinColumn({ name: 'channelId' })
  channel?: Channel;

  @Column()
  name: string;

  @OneToMany(() => ChannelCommand, _ => _.group)
  commands?: ChannelCommand[];

  @Column({ default: 'rgba(37, 38, 43, 1)' })
  color: string;
}