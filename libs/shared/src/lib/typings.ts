import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { HelixStream, HelixUser } from '@twurple/api';
import { rawDataSymbol } from '@twurple/common';

export type CachedStream = HelixStream[typeof rawDataSymbol] & { parsedMessages?: number };

export type AuthUser = HelixUser[typeof rawDataSymbol] & {
  dashboards: Array<DashboardAccess & { twitchUser: HelixUser[typeof rawDataSymbol] }>;
  isTester: boolean;
  isBotAdmin?: boolean;
};
