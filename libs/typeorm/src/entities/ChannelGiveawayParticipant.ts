import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { ChannelGiveaway } from './ChannelGiveaway';
// eslint-disable-next-line import/no-cycle
import { User } from './User';

@Entity('channels_giveaways_participants')
export class ChannelGiveawayParticipant {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('uuid')
  giveaway_id: string;

  @ManyToOne(() => ChannelGiveaway, (_) => _.participants)
  @JoinColumn({ name: 'giveaway_id' })
  giveaway?: ChannelGiveaway;

  @Column({ default: false })
  is_winner: boolean;

  @Column()
  user_id: string;

  @Column({ default: false })
  is_subscriber: boolean;

  @Column({ default: 1 })
  subscriber_tier: number;

  @ManyToOne(() => User, (_) => _.giveaway_participants)
  @JoinColumn({ name: 'user_id' })
  user?: User;

  @Column('timestamp', { nullable: true })
  user_follow_since: Date;

  @Column('bigint')
  user_stats_watched_time: bigint;

  @Column('integer', { name: 'messages', default: 0 })
  user_stats_messages: number;
}
