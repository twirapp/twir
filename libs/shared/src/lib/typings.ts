import { DashboardAccess } from '@tsuwari/typeorm/entities/DashboardAccess';
import { HelixStream, HelixUser } from '@twurple/api';
import { rawDataSymbol } from '@twurple/common';

export type CachedStream = HelixStream[typeof rawDataSymbol] & { parsedMessages?: number };

export type Dashboard = DashboardAccess & { twitch: HelixUser[typeof rawDataSymbol] };

export type AuthUser = HelixUser[typeof rawDataSymbol] & {
  dashboards: Dashboard[];
  isTester: boolean;
  isBotAdmin?: boolean;
};
