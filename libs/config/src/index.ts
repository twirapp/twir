import { resolve } from 'path';

import dotenv from 'dotenv';
import { cleanEnv, str, bool } from 'envalid';

try {
  dotenv.config({ path: resolve(process.cwd(), '.env') });
  // eslint-disable-next-line no-empty
} catch { }

export const config = cleanEnv(process.env, {
  DATABASE_URL: str({ default: 'postgresql://tsuwari:tsuwari@localhost:5432/tsuwari?schema=public' }),
  NODE_ENV: str({ choices: ['development', 'production'], default: 'development' }),
  TWITCH_CLIENTID: str(),
  TWITCH_CLIENTSECRET: str(),
  TWITCH_CALLBACKURL: str({ default: 'http://localhost:3005/login' }),
  JWT_EXPIRES_IN: str({ default: '5m' }),
  JWT_ACCESS_SECRET: str({ default: 'CoolSecretForAccess' }),
  JWT_REFRESH_SECRET: str({ default: 'CoolSecretForRefresh' }),
  REDIS_URL: str({ default: 'redis://:576294Aa@localhost:6379/0' }),
  SAY_IN_CHAT: bool({ default: true }),
  NATS_URL: str({ default: 'nats://localhost:4222' }),
});
