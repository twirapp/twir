import { SQL } from 'bun'

const MAX_STORAGE_SIZE = 30 * 1024 * 1024 // 30MB

export interface StorageOperation {
	action: 'get' | 'set' | 'delete' | 'has' | 'keys' | 'clear' | 'getTotalSize' |
		'push' | 'pop' | 'find' | 'filter' | 'splice' |
		'getProperty' | 'setProperty' | 'deleteProperty' | 'hasProperty' | 'merge'
	key?: string
	value?: unknown
	path?: string
	items?: unknown[]
	start?: number
	deleteCount?: number
	predicate?: string
}

export interface StorageResult {
	success: boolean
	data?: unknown
	error?: string
}

const sql = new SQL(process.env.DATABASE_URL!, { prepare: true })

export async function handleStorageOperation(
	channelId: string,
	op: StorageOperation
): Promise<StorageResult> {
	try {
		switch (op.action) {
			case 'get': {
				const rows = await sql`
					SELECT value FROM channels_storage
					WHERE channel_id = ${channelId} AND key = ${op.key}
				`
				if (rows.length === 0) return { success: true, data: null }
				const raw = rows[0].value
				const data = typeof raw === 'string' ? JSON.parse(raw) : raw
				return { success: true, data }
			}

			case 'set': {
				const valueJson = JSON.stringify(op.value)
				const currentSize = await getTotalSize(channelId)

				const existingRows = await sql`
					SELECT pg_column_size(value) as size FROM channels_storage
					WHERE channel_id = ${channelId} AND key = ${op.key}
				`
				const existingSize = existingRows.length > 0 ? (existingRows[0].size as number) : 0

				if (currentSize - existingSize + valueJson.length > MAX_STORAGE_SIZE) {
					return { success: false, error: 'Storage limit exceeded (30MB per channel)' }
				}

				await sql`
					INSERT INTO channels_storage (channel_id, key, value)
					VALUES (${channelId}, ${op.key!}, ${valueJson}::jsonb)
					ON CONFLICT (channel_id, key) DO UPDATE SET value = ${valueJson}::jsonb, updated_at = now()
				`
				return { success: true }
			}

			case 'delete': {
				const tag = await sql`
					DELETE FROM channels_storage
					WHERE channel_id = ${channelId} AND key = ${op.key}
				`
				return { success: true, data: tag.count > 0 }
			}

			case 'has': {
				const rows = await sql`
					SELECT 1 FROM channels_storage
					WHERE channel_id = ${channelId} AND key = ${op.key}
					LIMIT 1
				`
				return { success: true, data: rows.length > 0 }
			}

			case 'keys': {
				const rows = await sql`
					SELECT key FROM channels_storage
					WHERE channel_id = ${channelId}
					ORDER BY key
				`
				return { success: true, data: rows.map((r: any) => r.key) }
			}

			case 'clear': {
				await sql`
					DELETE FROM channels_storage
					WHERE channel_id = ${channelId}
				`
				return { success: true }
			}

			case 'getTotalSize': {
				const size = await getTotalSize(channelId)
				return { success: true, data: size }
			}

			default:
				return { success: false, error: `Unknown action: ${op.action}` }
		}
	} catch (err: any) {
		return { success: false, error: err.message || String(err) }
	}
}

async function getTotalSize(channelId: string): Promise<number> {
	const rows = await sql`
		SELECT COALESCE(SUM(pg_column_size(value)), 0) as total
		FROM channels_storage
		WHERE channel_id = ${channelId}
	`
	return Number(rows[0].total)
}
