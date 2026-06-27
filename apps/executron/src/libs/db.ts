import { createDecipheriv } from 'node:crypto'

import { config } from '@twir/config'
import { SQL } from 'bun'

const ALGORITHM = 'aes-256-gcm'
const NONCE_LENGTH = 12
const TAG_LENGTH = 16

const sql = new SQL(config.DATABASE_URL!, { prepare: true })

try {
	await sql`SELECT 1`
	console.log('Connected to database')
} catch (e) {
	console.error(e)
	process.exit(1)
}

const encryptionKey = Buffer.from(config.SECRETS_ENCRYPTION_KEY || '', 'utf-8')
if (encryptionKey.length !== 32) {
	console.error('SECRETS_ENCRYPTION_KEY must be 32 bytes')
	process.exit(1)
}

export async function getSecretsForChannel(channelId: string): Promise<Map<string, string>> {
	const result = await sql`
		SELECT name, value
		FROM channels_secrets
		WHERE "channel_id" = ${channelId}
	`

	const secrets = new Map<string, string>()

	for (const row of result) {
		try {
			const decrypted = decrypt(row.value as string)
			secrets.set(row.name as string, decrypted)
		} catch (error) {
			console.error(`Failed to decrypt secret "${row.name}":`, error)
		}
	}

	return secrets
}

function decrypt(encryptedBase64: string): string {
	const encryptedBuffer = Buffer.from(encryptedBase64, 'base64')

	const nonce = encryptedBuffer.subarray(0, NONCE_LENGTH)
	const authTag = encryptedBuffer.subarray(
		encryptedBuffer.length - TAG_LENGTH,
		encryptedBuffer.length
	)
	const ciphertext = encryptedBuffer.subarray(NONCE_LENGTH, encryptedBuffer.length - TAG_LENGTH)

	const decipher = createDecipheriv(ALGORITHM, encryptionKey, nonce)
	decipher.setAuthTag(authTag)

	let decrypted = decipher.update(ciphertext)
	decrypted = Buffer.concat([decrypted, decipher.final()])

	return decrypted.toString('utf-8')
}
