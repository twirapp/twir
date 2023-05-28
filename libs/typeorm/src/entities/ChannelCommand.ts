/* eslint-disable import/no-cycle */
import {
	Column,
	Entity,
	Index,
	JoinColumn,
	ManyToOne,
	OneToMany,
	PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';
import { ChannelCommandGroup } from './ChannelCommandGroup';
import { CommandResponse } from './CommandResponse';
import { CommandUsage } from './CommandUsage';

export enum CooldownType {
	GLOBAL = 'GLOBAL',
	PER_USER = 'PER_USER',
}

export enum CommandModule {
	CUSTOM = 'CUSTOM',
	DOTA = 'DOTA',
	MODERATION = 'MODERATION',
	MANAGE = 'MANAGE',
	SONGS = 'SONGS',
	TTS = 'TTS',
	STATS = 'STATS',
}

@Index('channels_commands_name_channelId_key', ['channelId', 'name'], { unique: true })
@Entity('channels_commands', { schema: 'public' })
export class ChannelCommand {
	@PrimaryGeneratedColumn('uuid')
	id: string;

	@Index()
	@Column('text', { name: 'name' })
	name: string;

	@Column('integer', { name: 'cooldown', nullable: true, default: 0 })
	cooldown: number | null;

	@Column('enum', { name: 'cooldownType', enum: CooldownType, default: CooldownType.GLOBAL })
	cooldownType: CooldownType;

	@Column('boolean', { name: 'enabled', default: true })
	enabled: boolean;

	@Column('text', { name: 'aliases', array: true, default: [] })
	aliases: string[];

	@Column('text', { name: 'description', nullable: true })
	description: string | null;

	@Column('boolean', { name: 'visible', default: true })
	visible: boolean;

	@Column('boolean', { name: 'is_reply', default: true })
	isReply: boolean;

	@Column('boolean', { name: 'default', default: false })
	default: boolean;

	@Column('text', { name: 'defaultName', nullable: true })
	defaultName: string | null;

	@Column('boolean', { default: true })
	keepResponsesOrder: boolean;

	@Column('enum', {
		name: 'module',
		enum: CommandModule,
		default: CommandModule.CUSTOM,
	})
	module: CommandModule;

	@ManyToOne(() => Channel, (_) => _.commands, {
		onDelete: 'RESTRICT',
		onUpdate: 'CASCADE',
	})
	@JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
	channel?: Channel;

	@Index()
	@Column('text', { name: 'channelId' })
	channelId: string;

	@OneToMany(() => CommandResponse, (_) => _.command)
	responses?: CommandResponse[];

	@OneToMany(() => CommandUsage, (_) => _.command)
	usages?: CommandUsage[];

	@ManyToOne(() => ChannelCommandGroup, (_) => _.commands, { onDelete: 'SET NULL' })
	group?: ChannelCommandGroup;

	@Column('uuid', { nullable: true })
	groupId?: string;

	@Column('text', { array: true, default: [] })
	deniedUsersIds: string[];

	@Column('text', { array: true, default: [] })
	allowedUsersIds: string[];

	@Column('text', { array: true, default: [] })
	rolesIds: string[];

	@Column('boolean', { default: false, name: 'online_only' })
	onlineOnly: boolean;

	@Column('int', { default: 0 })
	requiredWatchTime: number;

	@Column('int', { default: 0 })
	requiredMessages: number;

	@Column('int', { default: 0 })
	requiredUsedChannelPoints: number;
}
