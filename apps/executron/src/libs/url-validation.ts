import { lookup } from 'node:dns/promises'

const PRIVATE_RANGES: [number, number][] = [
	[0x7F000000, 0x7FFFFFFF],
	[0x0A000000, 0x0AFFFFFF],
	[0xAC100000, 0xAC1FFFFF],
	[0xC0A80000, 0xC0A8FFFF],
	[0xA9FE0000, 0xA9FEFFFF],
	[0x00000000, 0x00FFFFFF],
	[0xE0000000, 0xFFFFFFFF],
	[0x64400000, 0x647FFFFF],
]

const BLOCKED_HOSTNAMES = [
	'localhost',
	'metadata.google.internal',
	'metadata.goog',
]

function ipToLong(ip: string): number {
	return ip.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet, 10), 0) >>> 0
}

function isPrivateIp(ip: string): boolean {
	const long = ipToLong(ip)
	return PRIVATE_RANGES.some((range) => long >= range[0] && long <= range[1])
}

export async function validateUrl(rawUrl: string): Promise<void> {
	const parsed = new URL(rawUrl)

	if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
		throw new Error('Only http and https protocols are allowed')
	}

	const hostname = parsed.hostname

	if (BLOCKED_HOSTNAMES.includes(hostname)) {
		throw new Error(`Blocked hostname: ${hostname}`)
	}

	if (hostname.startsWith('twir_')) {
		throw new Error('Requests to internal services are not allowed')
	}

	if (/^\d+\.\d+\.\d+\.\d+$/.test(hostname)) {
		if (isPrivateIp(hostname)) {
			throw new Error(`Blocked private IP: ${hostname}`)
		}
		return
	}

	try {
		const { address } = await lookup(hostname, { family: 4 })
		if (isPrivateIp(address)) {
			throw new Error(`Blocked private IP: ${address}`)
		}
	} catch (err: any) {
		if (err.message?.startsWith('Blocked')) throw err
	}
}
