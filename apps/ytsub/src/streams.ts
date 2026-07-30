import { config } from '@twir/config'
import { SQL } from 'bun'

import type { ChannelBinding } from './message.ts'

const sql = new SQL(config.DATABASE_URL, {
	prepare: true,
	max: 20,
	idleTimeout: 30,
	maxLifetime: 3600,
	connectionTimeout: 10,
})

export interface YoutubeStream {
	readonly videoId: string
	readonly broadcasterName: string
	readonly title: string
	readonly viewers: number
	readonly startedAt: Date
}

export interface StreamBus {
	readonly Events: {
		readonly StreamOnline: {
			publish(message: {
				readonly startedAt: string
				readonly channelId: string
				readonly streamId: string
				readonly categoryName: string
				readonly categoryId: string
				readonly title: string
				readonly viewers: number
			}): Promise<void>
		}
		readonly StreamOffline: {
			publish(message: {
				readonly startedAt: string
				readonly channelId: string
			}): Promise<void>
		}
	}
}

export async function markOnline(bus: StreamBus, binding: ChannelBinding, stream: YoutubeStream): Promise<void> {
	await sql`
		INSERT INTO channels_streams (
			id, channel_id, "userId", "userLogin", "userName", "gameId", "gameName",
			"communityIds", type, title, "viewerCount", "startedAt", language,
			"thumbnailUrl", "tagIds", tags, "isMature", platform
		)
		VALUES (
			${stream.videoId}, ${binding.channelId}::uuid, ${binding.platformChannelId}, ${stream.broadcasterName},
			${stream.broadcasterName}, '', '', '{}', 'live', ${stream.title}, ${stream.viewers},
			${stream.startedAt}, '', '', '{}', '{}', FALSE, 'youtube'
		)
		ON CONFLICT ("userId", platform) DO UPDATE
		SET
			id = EXCLUDED.id,
			channel_id = EXCLUDED.channel_id,
			"userLogin" = EXCLUDED."userLogin",
			"userName" = EXCLUDED."userName",
			title = EXCLUDED.title,
			"viewerCount" = EXCLUDED."viewerCount",
			"startedAt" = EXCLUDED."startedAt"
	`

	await bus.Events.StreamOnline.publish({
		startedAt: stream.startedAt.toISOString(),
		channelId: binding.platformChannelId,
		streamId: stream.videoId,
		categoryName: '',
		categoryId: '',
		title: stream.title,
		viewers: stream.viewers,
	})
}

export async function markOffline(bus: StreamBus, platformChannelId: string, startedAt: Date): Promise<void> {
	await sql`
		DELETE FROM channels_streams
		WHERE "userId" = ${platformChannelId} AND platform = 'youtube'
	`

	await bus.Events.StreamOffline.publish({
		startedAt: startedAt.toISOString(),
		channelId: platformChannelId,
	})
}

export async function closeStreamsDatabase(): Promise<void> {
	await sql.close()
}
