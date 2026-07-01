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
	'169.254.169.254',
	'fd00:ec2::254',
	'instance-data',
]

const BLOCKED_DOMAIN_SUFFIXES = [
	'.internal',
	'.local',
	'.localhost',
	'.cluster.local',
	'.svc',
	'.default',
	'.kube',
	'.docker',
]

function ipToLong(ip: string): number {
	return ip.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet, 10), 0) >>> 0
}

function isPrivateIp(ip: string): boolean {
	const long = ipToLong(ip)
	return PRIVATE_RANGES.some((range) => long >= range[0] && long <= range[1])
}

function isPrivateIpv6(ip: string): boolean {
	const lower = ip.toLowerCase()
	if (lower === '::1' || lower === '0:0:0:0:0:0:0:1') return true
	if (lower.startsWith('fc') || lower.startsWith('fd')) return true
	if (lower.startsWith('fe80')) return true
	if (lower === '::' || lower === '0:0:0:0:0:0:0:0') return true
	if (lower.startsWith('::ffff:')) {
		const ipv4 = lower.slice(7)
		if (isPrivateIp(ipv4)) return true
	}
	return false
}

function isBlockedHostname(hostname: string): boolean {
	if (BLOCKED_HOSTNAMES.includes(hostname)) return true
	if (hostname.startsWith('twir_') || hostname.startsWith('twir-')) return true
	const lower = hostname.toLowerCase()
	return BLOCKED_DOMAIN_SUFFIXES.some((suffix) => lower.endsWith(suffix))
}

export interface ValidationResult {
	resolvedIp: string | null
}

export async function validateUrl(rawUrl: string): Promise<ValidationResult> {
	const parsed = new URL(rawUrl)

	if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
		throw new Error('Only http and https protocols are allowed')
	}

	const hostname = parsed.hostname

	if (isBlockedHostname(hostname)) {
		throw new Error(`Blocked hostname: ${hostname}`)
	}

	const isIpv6 = hostname.startsWith('[') && hostname.endsWith(']')
	if (isIpv6) {
		const bare = hostname.slice(1, -1)
		if (isPrivateIpv6(bare)) {
			throw new Error(`Blocked private IPv6: ${bare}`)
		}
		return { resolvedIp: null }
	}

	if (/^\d+\.\d+\.\d+\.\d+$/.test(hostname)) {
		if (isPrivateIp(hostname)) {
			throw new Error(`Blocked private IP: ${hostname}`)
		}
		return { resolvedIp: hostname }
	}

	try {
		const { address } = await lookup(hostname, { family: 4 })
		if (isPrivateIp(address)) {
			throw new Error(`Blocked private IP: ${address}`)
		}
		return { resolvedIp: address }
	} catch (err: any) {
		if (err.message?.startsWith('Blocked')) throw err
	}

	return { resolvedIp: null }
}
