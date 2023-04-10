import { resolve } from 'path';

import * as dotenv from 'dotenv';
import { bool, cleanEnv, str } from 'envalid';

try {
  dotenv.config({ path: resolve(process.cwd(), '../../.env') });
  // eslint-disable-next-line no-empty
} catch {}

export const config = cleanEnv(process.env, {
  DATABASE_URL: str({
    default: 'postgresql://tsuwari:tsuwari@postgres:5432/tsuwari?schema=public',
  }),
  NODE_ENV: str({ choices: ['development', 'production'], default: 'development' }),
  TWITCH_CLIENTID: str(),
  TWITCH_CLIENTSECRET: str(),
  TWITCH_CALLBACKURL: str({ default: 'http://localhost:3005/login' }),
  JWT_EXPIRES_IN: str({ default: '5m' }),
  JWT_ACCESS_SECRET: str({ default: 'CoolSecretForAccess' }),
  JWT_REFRESH_SECRET: str({ default: 'Cool`SecretForRefresh' }),
  REDIS_URL: str({ default: 'redis://localhost:6379/0' }),
  SAY_IN_CHAT: bool({ default: true }),
  HOSTNAME: str({ default: '' }),
  STEAM_USERNAME: str({ default: '' }),
  STEAM_PASSWORD: str({ default: '' }),
  STEAM_API_KEY: str({ default: '' }),
  MINIO_USER: str({ default: '' }),
  MINIO_PASSWORD: str({ default: '' }),
  MINIO_URL: str({ default: '' }),
  TOKENS_CIPHER_KEY: str({ default: 'pnyfwfiulmnqlhkvixaeligpprcnlyke' }),
  EVENTSUB_SECRET: str({ default: 'coolEventsubSecret' }),
  TTS_SERVICE_URL: str({ default: 'http://localhost:7000' }),
  SPOTIFY_CLIENT_ID: str({ default: '' }),
  SPOTIFY_CLIENT_SECRET: str({ default: '' }),
});
