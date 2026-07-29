import { config } from '@twir/config'
import { SQL } from 'bun'

import type { ChannelBinding, YoutubeTextChatMessage } from './message.ts'

const sql = new SQL(config.DATABASE_URL, {
	prepare: true,
	max: 20,
	idleTimeout: 30,
	maxLifetime: 3600,
	connectionTimeout: 10,
})

interface BindingRow {
	readonly id: string
	readonly channelId: string
	readonly platformChannelId: string
	readonly userId: string
	readonly enabled: boolean
}

interface UserRow {
	readonly id: string
}

export async function getYoutubeBinding(channelId: string): Promise<ChannelBinding | null> {
	const [binding] = await sql<BindingRow[]>`
		SELECT
			id::text AS "id",
			channel_id::text AS "channelId",
			platform_channel_id AS "platformChannelId",
			user_id::text AS "userId",
			enabled
		FROM channel_platforms
		WHERE channel_id = ${channelId}::uuid
			AND platform = 'youtube'
		LIMIT 1
	`

	return binding ?? null
}

export async function listYoutubeBindings(): Promise<readonly ChannelBinding[]> {
	return sql<BindingRow[]>`
		SELECT
			id::text AS "id",
			channel_id::text AS "channelId",
			platform_channel_id AS "platformChannelId",
			user_id::text AS "userId",
			enabled
		FROM channel_platforms
		WHERE platform = 'youtube' AND enabled = TRUE
	`
}

export async function ensureYoutubeChatter(
	binding: ChannelBinding,
	message: YoutubeTextChatMessage
): Promise<ChannelBinding> {
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

	await sql`
		INSERT INTO users_stats (user_id, channel_id, messages, is_mod, is_vip, is_subscriber)
		VALUES (
			${user.id}::uuid,
			${binding.channelId}::uuid,
			1,
			${message.author.isModerator},
			FALSE,
			FALSE
		)
		ON CONFLICT (user_id, channel_id) DO UPDATE
		SET
			messages = users_stats.messages + 1,
			is_mod = EXCLUDED.is_mod,
			is_vip = EXCLUDED.is_vip,
			is_subscriber = EXCLUDED.is_subscriber
	`

	return { ...binding, userId: user.id }
}

export async function closeDatabase(): Promise<void> {
	await sql.close()
}
