import * as crypto from 'node:crypto'

const ALGORITHM = 'aes-256-cbc'
const BLOCK_SIZE = 16

export function decrypt(text: string, key: string) {
	const contents = Buffer.from(text, 'hex')
	const iv = contents.slice(0, BLOCK_SIZE)
	const textBytes = contents.slice(BLOCK_SIZE)

	// @ts-expect-error TypeScript bug?
	const decipher = crypto.createDecipheriv(ALGORITHM, key, iv)
	let decrypted = decipher.update(textBytes.toString(), 'hex', 'utf8')
	decrypted += decipher.final('utf8')
	return decrypted
}

// Encrypts plain text into cipher text
export function encrypt(plainText: string, key: string) {
	const iv = crypto.randomBytes(BLOCK_SIZE)
	// @ts-expect-error TypeScript bug?
	const cipher = crypto.createCipheriv(ALGORITHM, key, iv)
	let cipherText
	try {
		cipherText = cipher.update(plainText, 'utf8', 'hex')
		cipherText += cipher.final('hex')
		cipherText = iv.toString('hex') + cipherText
	} catch (e) {
		console.log(e)
		cipherText = null
	}
	return cipherText
}
