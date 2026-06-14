import * as crypto from 'node:crypto'

const ALGORITHM = 'aes-256-cbc'
const BLOCK_SIZE = 16

export function decrypt(text: string, key: string) {
	const contents = Buffer.from(text, 'hex')
	const iv = new Uint8Array(contents.slice(0, BLOCK_SIZE))
	const textBytes = new Uint8Array(contents.slice(BLOCK_SIZE))

	const decipher = crypto.createDecipheriv(ALGORITHM, new Uint8Array(Buffer.from(key, 'utf8')), iv)
	let decrypted = decipher.update(Buffer.from(textBytes).toString('hex'), 'hex', 'utf8')
	decrypted += decipher.final('utf8')
	return decrypted
}

// Encrypts plain text into cipher text
export function encrypt(plainText: string, key: string) {
	const iv = new Uint8Array(crypto.randomBytes(BLOCK_SIZE))
	const cipher = crypto.createCipheriv(ALGORITHM, new Uint8Array(Buffer.from(key, 'utf8')), iv)
	let cipherText
	try {
		cipherText = cipher.update(plainText, 'utf8', 'hex')
		cipherText += cipher.final('hex')
		cipherText = Buffer.from(iv).toString('hex') + cipherText
	} catch (e) {
		console.log(e)
		cipherText = null
	}
	return cipherText
}
