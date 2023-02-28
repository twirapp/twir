import { HelixUserData } from '@twurple/api';

export type AuthUser = HelixUserData & {
  isBotAdmin?: boolean;
  apiKey: string;
};
