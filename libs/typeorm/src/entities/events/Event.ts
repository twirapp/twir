import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from 'typeorm';

// eslint-disable-next-line import/no-cycle
import { Channel } from '../Channel';
// eslint-disable-next-line import/no-cycle
import { EventOperation } from './EventOperation';

export enum EventType {
	FOLLOW = 'FOLLOW',
	SUBSCRIBE = 'SUBSCRIBE',
	RESUBSCRIBE = 'RESUBSCRIBE',
	SUB_GIFT = 'SUB_GIFT',
	REDEMPTION_CREATED = 'REDEMPTION_CREATED',
	COMMAND_USED = 'COMMAND_USED',
	FIRST_USER_MESSAGE = 'FIRST_USER_MESSAGE',
	RAIDED = 'RAIDED',
	TITLE_OR_CATEGORY_CHANGED = 'TITLE_OR_CATEGORY_CHANGED',
	STREAM_ONLINE = 'STREAM_ONLINE',
	STREAM_OFFLINE = 'STREAM_OFFLINE',
	ON_CHAT_CLEAR = 'ON_CHAT_CLEAR',
	DONATE = 'DONATE',
	KEYWORD_MATCHED = 'KEYWORD_MATCHED',
	GREETING_SENDED = 'GREETING_SENDED',
	POLL_BEGIN = 'POLL_BEGIN',
	POLL_PROGRESS = 'POLL_PROGRESS',
	POLL_END = 'POLL_END',
	PREDICTION_BEGIN = 'PREDICTION_BEGIN',
	PREDICTION_PROGRESS = 'PREDICTION_PROGRESS',
	PREDICTION_END = 'PREDICTION_END',
	PREDICTION_LOCK = 'PREDICTION_LOCK',
	STREAM_FIRST_USER_JOIN = 'STREAM_FIRST_USER_JOIN',
}

@Entity({ name: 'channels_events' })
export class Event {
	@PrimaryGeneratedColumn('uuid')
	id: string;

	@Column('enum', { enum: EventType })
	type: EventType;

	@Column('text', { nullable: true })
	description: string | null;

	@Column('uuid', { nullable: true })
	rewardId: string | null;

	@Column('text', { nullable: true })
	commandId: string | null;

	@Column('text', { nullable: true })
	keywordId: string | null;

	@OneToMany(() => EventOperation, (_) => _.event)
	operations: EventOperation[];

	@ManyToOne(() => Channel, (_) => _.events)
	@JoinColumn({ name: 'channelId' })
	channel?: Channel;

	@Column()
	channelId: string;

	@Column({ default: true })
	enabled: boolean;

	@Column('boolean', { default: false, name: 'online_only' })
	onlineOnly: boolean;
}
