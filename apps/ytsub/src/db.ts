import { config } from '@twir/config'
import { SQL } from 'bun'

import { isYoutubeSubscriberBadge } from './message.ts'

import type { ChannelBinding, YoutubeTextChatMessage } from './message.ts'

const sql = new SQL(config.DATABASE_URL, {
	prepare: true,
	max: 20,
	idleTimeout: 30,
	maxLifetime: 3600,
	connectionTimeout: 10,
})

const CHATTER_CACHE_MAX_SIZE = 10_000

interface CachedChatter {
	readonly internalUserId: string
	readonly name: string
}

const chatterCache = new Map<string, CachedChatter>()

interface BindingRow {
	readonly id: string
	readonly channelId: string
	readonly platformChannelId: string
	readonly botPlatformId: string | null
	readonly userId: string
	readonly enabled: boolean
}

interface UserRow {
	readonly id: string
}

export async function getYoutubeBinding(channelId: string): Promise<ChannelBinding | null> {
	const [binding] = await sql<BindingRow[]>`
		SELECT
			channel_platforms.id::text AS "id",
			channel_platforms.channel_id::text AS "channelId",
			channel_platforms.platform_channel_id AS "platformChannelId",
			bot_user.platform_id AS "botPlatformId",
			channel_platforms.user_id::text AS "userId",
			channel_platforms.enabled
		FROM channel_platforms
		LEFT JOIN users AS bot_user
			ON bot_user.id = channel_platforms.bot_user_id
			AND bot_user.platform = 'youtube'
		WHERE channel_platforms.channel_id = ${channelId}::uuid
			AND channel_platforms.platform = 'youtube'
		LIMIT 1
	`

	return binding ?? null
}

export async function listYoutubeBindings(): Promise<readonly ChannelBinding[]> {
	return sql<BindingRow[]>`
		SELECT
			channel_platforms.id::text AS "id",
			channel_platforms.channel_id::text AS "channelId",
			channel_platforms.platform_channel_id AS "platformChannelId",
			bot_user.platform_id AS "botPlatformId",
			channel_platforms.user_id::text AS "userId",
			channel_platforms.enabled
		FROM channel_platforms
		LEFT JOIN users AS bot_user
			ON bot_user.id = channel_platforms.bot_user_id
			AND bot_user.platform = 'youtube'
		WHERE channel_platforms.platform = 'youtube' AND channel_platforms.enabled = TRUE
	`
}

export async function ensureYoutubeChatter(
	binding: ChannelBinding,
	message: YoutubeTextChatMessage
): Promise<ChannelBinding> {
	const cached = chatterCache.get(message.author.id)
	let internalUserId: string
	if (cached?.name === message.author.name) {
		internalUserId = cached.internalUserId
		chatterCache.delete(message.author.id)
		chatterCache.set(message.author.id, cached)
	} else {
		const [user] = await sql<UserRow[]>`
			INSERT INTO users (platform, platform_id, login, display_name)
			VALUES ('youtube', ${message.author.id}, ${message.author.name}, ${message.author.name})
			ON CONFLICT (platform, platform_id) DO UPDATE
			SET login = EXCLUDED.login, display_name = EXCLUDED.display_name
			RETURNING id::text AS "id"
		`
		if (!user) {
			throw new Error('YouTube user upsert returned no row')
		}
		internalUserId = user.id
		if (chatterCache.size >= CHATTER_CACHE_MAX_SIZE) {
			const oldestPlatformId = chatterCache.keys().next().value
			if (oldestPlatformId !== undefined) {
				chatterCache.delete(oldestPlatformId)
			}
		}
		chatterCache.set(message.author.id, { internalUserId, name: message.author.name })
	}

	await sql`
		INSERT INTO users_stats (user_id, channel_id, messages, is_mod, is_vip, is_subscriber)
		VALUES (
			${internalUserId}::uuid,
			${binding.channelId}::uuid,
			1,
			${message.author.isModerator},
			FALSE,
			${message.author.badges.some(isYoutubeSubscriberBadge)}
		)
		ON CONFLICT (user_id, channel_id) DO UPDATE
		SET
			messages = users_stats.messages + 1,
			is_mod = EXCLUDED.is_mod,
			is_vip = EXCLUDED.is_vip,
			is_subscriber = EXCLUDED.is_subscriber
	`

	return { ...binding, userId: internalUserId }
}

export async function closeDatabase(): Promise<void> {
	await sql.close()
}
