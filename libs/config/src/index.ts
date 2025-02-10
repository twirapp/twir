import { z } from 'zod'

const envSchema = z.object({
	DATABASE_URL: z.string().default('postgresql://tsuwari:tsuwari@localhost:54321/tsuwari?schema=public'),
	NODE_ENV: z.enum(['development', 'production', 'test']).default('development'),
	TWITCH_CLIENTID: z.string(),
	TWITCH_CLIENTSECRET: z.string(),
	TWITCH_CALLBACKURL: z.string().default('http://localhost:3005/login'),
	REDIS_URL: z.string().default('redis://localhost:6379/0'),
	SAY_IN_CHAT: z.boolean().default(true),
	SITE_BASE_URL: z.string().default('localhost:3005'),
	STEAM_USERNAME: z.string().optional().default(undefined),
	STEAM_PASSWORD: z.string().default(undefined),
	STEAM_API_KEY: z.string().default(undefined),
	MINIO_USER: z.string().optional().default(undefined),
	MINIO_PASSWORD: z.string().optional().default(undefined),
	MINIO_URL: z.string().optional().default(undefined),
	TOKENS_CIPHER_KEY: z.string().default('pnyfwfiulmnqlhkvixaeligpprcnlyke'),
	EVENTSUB_SECRET: z.string().default('coolEventsubSecret'),
	TTS_SERVICE_URL: z.string().default('http://localhost:7001'),
	SPOTIFY_CLIENT_ID: z.string().optional().default(undefined),
	SPOTIFY_CLIENT_SECRET: z.string().optional().default(undefined),
	ODESLI_API_KEY: z.string().optional().default(undefined),
	DISCORD_FEEDBACK_URL: z.string().default(undefined),
	NATS_URL: z.string().default('127.0.0.1:4222'),
	USE_WSS: z.enum(['true', 'false']).transform((value) => value === 'true'),
})

export const config = envSchema.parse(Bun.env)
