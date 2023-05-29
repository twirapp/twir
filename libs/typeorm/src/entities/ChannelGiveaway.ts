import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Channel } from './Channel';
// eslint-disable-next-line import/no-cycle
import { ChannelGiveawayParticipant } from './ChannelGiveawayParticipant';

export enum ChannelGiveAwayType {
	BY_KEYWORD = 'BY_KEYWORD',
	BY_RANDOM_NUMBER = 'BY_RANDOM_NUMBER',
}

@Entity('channels_giveaways')
export class ChannelGiveaway {
	@PrimaryGeneratedColumn('uuid')
	id: string;

	@Column('text')
	description: string;

	@Column('enum', { enum: ChannelGiveAwayType })
	type: ChannelGiveAwayType;

	@Column()
	channel_id: string;

	@ManyToOne(() => Channel, (_) => _.giveAways)
	@JoinColumn({ name: 'channel_id' })
	channel?: Channel;

	@Column('timestamp')
	created_at: Date;

	@Column('timestamp')
	start_at: Date;

	@Column('timestamp', { nullable: true })
	end_at: Date | null;

	@Column('timestamp')
	closed_at: Date;

	@Column({ nullable: true })
	required_min_watch_time: number | null;

	@Column({ nullable: true })
	required_min_follow_time: number | null;

	@Column({ nullable: true })
	required_min_messages: number | null;

	@Column({ nullable: true })
	required_min_subscriber_tier: number | null;

	@Column({ nullable: true })
	required_min_subscribe_time: number | null;

	@Column('simple-array')
	eligible_user_groups: string[];

	@Column({ nullable: true })
	keyword: string | null;

	@Column({ nullable: true })
	random_number_from: number | null;

	@Column({ nullable: true })
	random_number_to: number | null;

	@Column({ nullable: true })
	winning_random_number: number | null;

	@Column()
	winners_count: number;

	@Column({ default: 0 })
	subscribers_luck: number;

	@Column({ default: 0 })
	subscribers_tier1_luck: number;

	@Column({ default: 0 })
	subscribers_tier2_luck: number;

	@Column({ default: 0 })
	subscribers_tier3_luck: number;

	@Column('simple-array', {
		default: [],
		array: true,
	})
	watched_time_lucks: ChannelGiveawayConfigurableLuck[];

	@Column('simple-array', {
		default: [],
		array: true,
	})
	messages_lucks: ChannelGiveawayConfigurableLuck[];

	@Column('simple-array', {
		default: [],
		array: true,
	})
	used_channel_points_lucks: ChannelGiveawayConfigurableLuck[];

	@OneToMany(() => ChannelGiveawayParticipant, (_) => _.giveaway)
	participants?: ChannelGiveawayParticipant[];
}

type ChannelGiveawayConfigurableLuck = {
	value: number;
	luck: number;
};
