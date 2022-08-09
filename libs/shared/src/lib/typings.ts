import { DashboardAccess } from '@tsuwari/prisma';
import { HelixStreamData, HelixUserData } from '@twurple/api';

export type CachedStream = HelixStreamData & { parsedMessages?: number };

export type AuthUser = HelixUserData & {
  dashboards: Array<DashboardAccess & { twitch: HelixUserData }>;
  isTester: boolean;
  isBotAdmin?: boolean;
};
