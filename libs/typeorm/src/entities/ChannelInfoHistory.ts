import { Column, CreateDateColumn, Entity, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Channel } from './Channel';

@Entity({ name: 'channels_info_history' })
export class ChannelInfoHistory {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  channelId: string;

  @ManyToOne(() => Channel, _ => _.infoHistories)
  channel?: Channel;

  @CreateDateColumn()
  createdAt: Date;

  @Column('text')
  title: string;

  @Column()
  category: string;
}