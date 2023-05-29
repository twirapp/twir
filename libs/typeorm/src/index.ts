import { resolve } from 'path';
import 'reflect-metadata';

import * as dotenv from 'dotenv';
import { DataSource } from 'typeorm';

import { Bot } from './entities/Bot';
import { Channel } from './entities/Channel';
import { ChannelChatMessage } from './entities/ChannelChatMessage';
import { ChannelCommand } from './entities/ChannelCommand';
import { ChannelCommandGroup } from './entities/ChannelCommandGroup';
import { ChannelCustomvar } from './entities/ChannelCustomvar';
import { ChannelDotaAccount } from './entities/ChannelDotaAccount';
import { ChannelEmoteUsage } from './entities/ChannelEmoteUsage';
import { ChannelEvent } from './entities/ChannelEvent';
import { ChannelDonationEvent } from './entities/channelEvents/Donation';
import { ChannelFollowEvent } from './entities/channelEvents/Follow';
import { ChannelGreeting } from './entities/ChannelGreeting';
import { ChannelInfoHistory } from './entities/ChannelInfoHistory';
import { ChannelIntegration } from './entities/ChannelIntegration';
import { ChannelKeyword } from './entities/ChannelKeyword';
import { ChannelModerationSetting } from './entities/ChannelModerationSetting';
import { ChannelModerationWarn } from './entities/ChannelModerationWarn';
import { ChannelModuleSettings } from './entities/ChannelModuleSettings';
import { ChannelPermit } from './entities/ChannelPermit';
import { ChannelRole } from './entities/ChannelRole';
import { ChannelRoleUser } from './entities/ChannelRoleUser';
import { ChannelStream } from './entities/ChannelStream';
import { ChannelTimer } from './entities/ChannelTimer';
import { ChannelTimerResponse } from './entities/ChannelTimerResponse';
import { CommandResponse } from './entities/CommandResponse';
import { CommandUsage } from './entities/CommandUsage';
import { DotaGameMode } from './entities/DotaGameMode';
import { DotaHero } from './entities/DotaHero';
import { DotaMatch } from './entities/DotaMatch';
import { DotaMatchCard } from './entities/DotaMatchCard';
import { DotaMatchResult } from './entities/DotaMatchResult';
import { Event } from './entities/events/Event';
import { EventOperation } from './entities/events/EventOperation';
import { EventOperationFilter } from './entities/events/EventOperationFilter';
import { IgnoredUser } from './entities/IgnoredUser';
import { Integration } from './entities/Integration';
import { Notification } from './entities/Notification';
import { NotificationMessage } from './entities/NotificationMessage';
import { RequestedSong } from './entities/RequestedSong';
import { Token } from './entities/Token';
import { User } from './entities/User';
import { UserFile } from './entities/UserFile';
import { UserOnline } from './entities/UserOnline';
import { UserStats } from './entities/UserStats';
import { UserViewedNotification } from './entities/UserViewedNotification';
import { ChannelGiveaway } from './entities/ChannelGiveaway';
import { ChannelGiveawayParticipant } from './entities/ChannelGiveawayParticipant';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

export * from 'typeorm';

export const AppDataSource = new DataSource({
	type: 'postgres',
	url: process.env.DATABASE_URL,
	entities: [
		Bot,
		Channel,
		ChannelCommand,
		ChannelCustomvar,
		ChannelDotaAccount,
		ChannelGreeting,
		ChannelIntegration,
		ChannelKeyword,
		ChannelModerationSetting,
		ChannelModerationWarn,
		ChannelModuleSettings,
		ChannelPermit,
		ChannelTimer,
		ChannelTimerResponse,
		ChannelGiveaway,
		ChannelGiveawayParticipant,
		CommandResponse,
		CommandUsage,
		DotaGameMode,
		DotaHero,
		DotaMatch,
		DotaMatchCard,
		DotaMatchResult,
		Integration,
		Notification,
		NotificationMessage,
		Token,
		User,
		UserFile,
		UserStats,
		UserViewedNotification,
		UserOnline,
		ChannelEvent,
		ChannelFollowEvent,
		ChannelDonationEvent,
		ChannelStream,
		ChannelChatMessage,
		RequestedSong,
		IgnoredUser,
		ChannelEmoteUsage,
		Event,
		EventOperation,
		ChannelInfoHistory,
		ChannelCommandGroup,
		ChannelRole,
		ChannelRoleUser,
		EventOperationFilter,
	],
	subscribers: [],
	migrations: ['src/migrations/*.ts'],
	migrationsTableName: 'typeorm_migrations',
});
