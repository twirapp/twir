import { resolve } from 'node:path'

import * as dotenv from 'dotenv'
import { bool, cleanEnv, str } from 'envalid'

try {
	dotenv.config({ path: resolve(process.cwd(), '../../.env') })
} catch {
}

export const config = cleanEnv(process.env, {
	DATABASE_URL: str({
		default: 'postgresql://tsuwari:tsuwari@localhost:54321/tsuwari?schema=public',
	}),
	NODE_ENV: str({ choices: ['development', 'production'], default: 'development' }),
	TWITCH_CLIENTID: str({ default: '' }),
	TWITCH_CLIENTSECRET: str({ default: '' }),
	TWITCH_CALLBACKURL: str({ default: 'http://localhost:3005/login' }),
	REDIS_URL: str({ default: 'redis://localhost:6379/0' }),
	SAY_IN_CHAT: bool({ default: true }),
	SITE_BASE_URL: str({ default: 'localhost:3005' }),
	STEAM_USERNAME: str({ default: '' }),
	STEAM_PASSWORD: str({ default: '' }),
	STEAM_API_KEY: str({ default: '' }),
	MINIO_USER: str({ default: '' }),
	MINIO_PASSWORD: str({ default: '' }),
	MINIO_URL: str({ default: '' }),
	TOKENS_CIPHER_KEY: str({ default: 'pnyfwfiulmnqlhkvixaeligpprcnlyke' }),
	EVENTSUB_SECRET: str({ default: 'coolEventsubSecret' }),
	TTS_SERVICE_URL: str({ default: 'http://localhost:7001' }),
	SPOTIFY_CLIENT_ID: str({ default: '' }),
	SPOTIFY_CLIENT_SECRET: str({ default: '' }),
	ODESLI_API_KEY: str({ default: '' }),
	DISCORD_FEEDBACK_URL: str({ default: '' }),
	NATS_URL: str({ default: '127.0.0.1:4222' }),
	USE_WSS: bool({ default: false }),
})
