import { resolve } from 'path';

import dotenv from 'dotenv';
import { cleanEnv, str, bool } from 'envalid';

try {
  dotenv.config({ path: resolve(process.cwd(), '.env') });
  // eslint-disable-next-line no-empty
} catch { }

export const config = cleanEnv(process.env, {
  NODE_ENV: str({ choices: ['development', 'production'], default: 'development' }),
  TWITCH_CLIENTID: str(),
  TWITCH_CLIENTSECRET: str(),
  TWITCH_CALLBACKURL: str({ default: 'http://localhost:3006/login' }),
  JWT_EXPIRES_IN: str({ default: '5m' }),
  JWT_ACCESS_SECRET: str({ default: 'CoolSecretForAccess' }),
  JWT_REFRESH_SECRET: str({ default: 'CoolSecretForRefresh' }),
  REDIS_URL: str({ default: 'redis://:576294Aa@localhost:6379/0' }),
  SAY_IN_CHAT: bool({ default: true }),
  MICROSERVICE_STREAM_STATUS_URL: str({ default: '0.0.0.0:50000' }),
  MICROSERVICE_BOTS_URL: str({ default: '0.0.0.0:50001' }),
  MICROSERVICE_WATCHED_URL: str({ default: '0.0.0.0:50002' }),
});
